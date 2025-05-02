package extrude

import (
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	// "log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-colours"
	"github.com/aaronland/go-colours/extruder"
	"github.com/aaronland/go-colours/grid"
	"github.com/aaronland/go-colours/palette"
)

type Closest struct {
	Palette string         `json:"palette"`
	Colour  colours.Colour `json:"colour"`
}

type Swatch struct {
	Colour  colours.Colour `json:"colour"`
	Closest []*Closest     `json:"closest"`
}

type Extrusion struct {
	Extruder string    `json:"extruder"`
	Palettes []string  `json:"palettes"`
	Swatches []*Swatch `json:"swatches"`
}

type Image struct {
	URI        string       `json:"uri"`
	Extrusions []*Extrusion `json:"extrusions"`
}

type ExtrudeOptions struct {
	// A list of aaronland/go-colours/extruder.Extruder URIs to instantiate.
	ExtruderURIs []string
	// A list of aaronland/go-colours/palette.Palette URIs to instantiate.
	PaletteURIs []string
	// The directory where temporary or derivative image files should be written. If empty a temporary directory will be created.
	Root string
	// The list of images to extrude colour palettes for.
	Images []string
	// Allow fetching remote images via HTTP or HTTPS.
	AllowRemote bool
	// Create derivative copies of images for post-processing operations, for example the HTML page used by the app/review code.
	CloneImages bool
}

type ExtrudeResponse struct {
	// The list of final `Image` instances for which colour palettes were derived.
	Images []*Image
	// The list of palette names for which "snap-to-grid" matches were derived.
	Palettes []string
	// The absolute path of the directory where derivative image files were written.
	Root string
	// A boolean flag indicating whether `Root` (the directory where derivative image files were written) was created at runtime as a temporary directory.
	IsTmpRoot bool
}

// Extrude (derive) dominant colours from one or more images as well as closest matches colours using zero or more "snap-to-grid" colour palettes.
func Extrude(ctx context.Context, opts *ExtrudeOptions) (*ExtrudeResponse, error) {

	var abs_root string
	var tmp_root bool

	if opts.Root != "" {

		root_dir, err := filepath.Abs(opts.Root)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive root, %w", err)
		}

		abs_root = root_dir
		tmp_root = false

	} else {

		root_dir, err := os.MkdirTemp("", "extrude")

		if err != nil {
			return nil, fmt.Errorf("Failed to create temp dir, %w", err)
		}

		abs_root = root_dir
		tmp_root = true
	}

	extruders := make([]extruder.Extruder, len(opts.ExtruderURIs))

	for idx, ex_uri := range opts.ExtruderURIs {

		ex, err := extruder.NewExtruder(ctx, ex_uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to create new '%s' extruder, %w", ex_uri, err)
		}

		extruders[idx] = ex
	}

	gr, err := grid.NewGrid(ctx, "euclidian://")

	if err != nil {
		return nil, fmt.Errorf("Failed to create new grid, %w", err)
	}

	palettes := make([]palette.Palette, len(opts.PaletteURIs))

	for idx, p_uri := range opts.PaletteURIs {

		p, err := palette.NewPalette(ctx, p_uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to create '%s' palette, %w", p_uri, err)
		}

		palettes[idx] = p
	}

	derive_colours := func(im image.Image) ([]*Extrusion, error) {

		extrusions := make([]*Extrusion, 0)

		for _, ex := range extruders {

			swatches := make([]*Swatch, 0)

			colours, err := ex.Colours(ctx, im, 5)

			if err != nil {
				return nil, fmt.Errorf("Failed to derive colours for image, %w", err)
			}

			for _, c := range colours {

				closest := make([]*Closest, 0)

				for _, p := range palettes {

					c2, err := gr.Closest(ctx, c, p)

					if err != nil {
						return nil, fmt.Errorf("Failed to derive closest, %w", err)
					}

					cl := &Closest{
						Palette: p.Reference(),
						Colour:  c2,
					}

					closest = append(closest, cl)
				}

				sw := &Swatch{
					Colour:  c,
					Closest: closest,
				}

				swatches = append(swatches, sw)
			}

			palette_labels := make([]string, 0)

			for _, p := range palettes {
				palette_labels = append(palette_labels, p.Reference())
			}

			e := &Extrusion{
				Extruder: ex.Name(),
				Palettes: palette_labels,
				Swatches: swatches,
			}

			extrusions = append(extrusions, e)
		}

		return extrusions, nil
	}

	images := make([]*Image, 0)

	for _, path := range opts.Images {

		if strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "http://") {

			if !opts.AllowRemote {
				return nil, fmt.Errorf("Remote images not allowed")
			}

			fname := filepath.Base(path)

			rsp, err := http.Get(path)

			if err != nil {
				return nil, fmt.Errorf("Failed to fetch %s, %w", path, err)
			}

			defer rsp.Body.Close()

			new_path := filepath.Join(abs_root, fname)
			new_wr, err := os.OpenFile(new_path, os.O_RDWR|os.O_CREATE, 0644)

			if err != nil {
				return nil, fmt.Errorf("Failed to open %s for writing, %w", new_path, err)
			}

			_, err = io.Copy(new_wr, rsp.Body)

			if err != nil {
				return nil, fmt.Errorf("Failed to copy %s to %s, %w", path, new_path, err)
			}

			err = new_wr.Close()

			if err != nil {
				return nil, fmt.Errorf("Failed to close %s after writing, %w", new_path, err)
			}

			path = new_path
		}

		fname := filepath.Base(path)
		ext := filepath.Ext(fname)
		fname = strings.Replace(fname, ext, ".png", 1)

		r, err := os.Open(path)

		if err != nil {
			return nil, fmt.Errorf("Failed to open %s for reading, %w", path, err)
		}

		defer r.Close()

		im, _, err := image.Decode(r)

		if err != nil {
			return nil, fmt.Errorf("Failed to decode %s, %w", path, err)
		}

		extrusions, err := derive_colours(im)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive colours, %w", err)
		}

		im_c := &Image{
			URI:        fname,
			Extrusions: extrusions,
		}

		images = append(images, im_c)

		if opts.CloneImages {

			im_path := filepath.Join(abs_root, fname)

			im_wr, err := os.OpenFile(im_path, os.O_RDWR|os.O_CREATE, 0644)

			if err != nil {
				return nil, fmt.Errorf("Failed to open %s for writing, %w", im_path, err)
			}

			err = png.Encode(im_wr, im)

			if err != nil {
				return nil, fmt.Errorf("Failed to encode %s, %w", im_path, err)
			}

			err = im_wr.Close()

			if err != nil {
				return nil, fmt.Errorf("Failed to close %s after writing, %w", im_path, err)
			}
		}
	}

	str_palettes := make([]string, 0)

	for _, p := range palettes {
		str_palettes = append(str_palettes, p.Reference())
	}

	rsp := &ExtrudeResponse{
		Images:    images,
		Palettes:  str_palettes,
		Root:      abs_root,
		IsTmpRoot: tmp_root,
	}

	return rsp, nil
}
