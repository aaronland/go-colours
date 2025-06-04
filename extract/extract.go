package colours

import (
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
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
	Palette string
	Colour  colours.Colour
}

type Swatch struct {
	Colour  colours.Colour
	Closest []*Closest
}

type Extrusion struct {
	Extruder string
	Palettes []string
	Swatches []*Swatch
}

type Image struct {
	URI        string
	Extrusions []*Extrusion
}

type ExtractOptions struct {
	ExtruderURIs []string
	PaletteURIs  []string
	Root         string
	Images       []string
	AllowRemote  bool
}

func Extract(ctx context.Context, opts *ExtractOptions) ([]*Image, error) {

	var abs_root string

	if opts.Root != "" {

		root_dir, err := filepath.Abs(opts.Root)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive root, %w", err)
		}

		abs_root = root_dir

	} else {

		root_dir, err := os.MkdirTemp("", "extract")

		if err != nil {
			return nil, fmt.Errorf("Failed to create temp dir, %w", err)
		}

		defer os.RemoveAll(root_dir)
		abs_root = root_dir
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

		/*
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
		*/
	}

	return images, nil
}
