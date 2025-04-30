// Command-line tool to generate an HTML page (and associated assets) to review the colour extraction
// for an image using one or more extruders and one or more palettes. The application will spawn a short-lived
// web server to serve the HTML review on a random port number and open its URI in the default browser.
package review

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-colours"
	"github.com/aaronland/go-colours/extruder"
	"github.com/aaronland/go-colours/grid"
	"github.com/aaronland/go-colours/palette"
	"github.com/sfomuseum/go-www-show"
)

//go:embed index.html
var index_html string

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

type TemplateVars struct {
	Images   []*Image
	Palettes []string
}

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(fs)

	if err != nil {
		return err
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	var abs_root string

	if opts.Root != "" {

		root_dir, err := filepath.Abs(opts.Root)

		if err != nil {
			return fmt.Errorf("Failed to derive root, %w", err)
		}

		abs_root = root_dir

	} else {

		root_dir, err := os.MkdirTemp("", "colours")

		if err != nil {
			return fmt.Errorf("Failed to create temp dir, %w", err)
		}

		defer os.RemoveAll(root_dir)
		abs_root = root_dir
	}

	extruders := make([]extruder.Extruder, len(opts.ExtruderURIs))

	for idx, ex_uri := range opts.ExtruderURIs {

		ex, err := extruder.NewExtruder(ctx, ex_uri)

		if err != nil {
			return fmt.Errorf("Failed to create new '%s' extruder, %w", ex_uri, err)
		}

		extruders[idx] = ex
	}

	gr, err := grid.NewGrid(ctx, "euclidian://")

	if err != nil {
		return fmt.Errorf("Failed to create new grid, %w", err)
	}

	palettes := make([]palette.Palette, len(opts.PaletteURIs))

	for idx, p_uri := range opts.PaletteURIs {

		p, err := palette.NewPalette(ctx, p_uri)

		if err != nil {
			return fmt.Errorf("Failed to create '%s' palette, %w", p_uri, err)
		}

		palettes[idx] = p
	}

	index_t, err := template.New("index").Parse(index_html)

	if err != nil {
		return fmt.Errorf("Failed to parse template, %w", err)
	}

	derive_colours := func(im image.Image) ([]*Extrusion, error) {

		extrusions := make([]*Extrusion, 0)

		for _, ex := range extruders {

			swatches := make([]*Swatch, 0)

			colours, err := ex.Colours(im, 5)

			if err != nil {
				return nil, fmt.Errorf("Failed to derive colours for image, %w", err)
			}

			for _, c := range colours {

				closest := make([]*Closest, 0)

				for _, p := range palettes {

					c2, err := gr.Closest(c, p)

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

			fname := filepath.Base(path)

			rsp, err := http.Get(path)

			if err != nil {
				return fmt.Errorf("Failed to fetch %s, %w", path, err)
			}

			defer rsp.Body.Close()

			new_path := filepath.Join(abs_root, fname)
			new_wr, err := os.OpenFile(new_path, os.O_RDWR|os.O_CREATE, 0644)

			if err != nil {
				return fmt.Errorf("Failed to open %s for writing, %w", new_path, err)
			}

			_, err = io.Copy(new_wr, rsp.Body)

			if err != nil {
				return fmt.Errorf("Failed to copy %s to %s, %w", path, new_path, err)
			}

			err = new_wr.Close()

			if err != nil {
				return fmt.Errorf("Failed to close %s after writing, %w", new_path, err)
			}

			path = new_path
		}

		fname := filepath.Base(path)
		ext := filepath.Ext(fname)
		fname = strings.Replace(fname, ext, ".png", 1)

		r, err := os.Open(path)

		if err != nil {
			return fmt.Errorf("Failed to open %s for reading, %w", path, err)
		}

		defer r.Close()

		im, _, err := image.Decode(r)

		if err != nil {
			return fmt.Errorf("Failed to decode %s, %w", path, err)
		}

		extrusions, err := derive_colours(im)

		if err != nil {
			return fmt.Errorf("Failed to derive colours, %w", err)
		}

		im_c := &Image{
			URI:        fname,
			Extrusions: extrusions,
		}

		images = append(images, im_c)

		im_path := filepath.Join(abs_root, fname)

		im_wr, err := os.OpenFile(im_path, os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			return fmt.Errorf("Failed to open %s for writing, %w", im_path, err)
		}

		err = png.Encode(im_wr, im)

		if err != nil {
			return fmt.Errorf("Failed to encode %s, %w", im_path, err)
		}

		err = im_wr.Close()

		if err != nil {
			return fmt.Errorf("Failed to close %s after writing, %w", im_path, err)
		}

	}

	//

	index_path := filepath.Join(abs_root, "index.html")

	index_wr, err := os.OpenFile(index_path, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return fmt.Errorf("Failed to open %s for writing, %w", index_path, err)
	}

	str_palettes := make([]string, 0)

	for _, p := range palettes {
		str_palettes = append(str_palettes, p.Reference())
	}

	vars := TemplateVars{
		Images:   images,
		Palettes: str_palettes,
	}

	err = index_t.Execute(index_wr, vars)

	if err != nil {
		return fmt.Errorf("Failed to encode %s, %w", index_path, err)
	}

	err = index_wr.Close()

	if err != nil {
		return fmt.Errorf("Failed to close %s after writing, %w", index_path, err)
	}

	//

	mux := http.NewServeMux()

	dir_fs := os.DirFS(abs_root)
	http_fs := http.FileServerFS(dir_fs)

	mux.Handle("/", http_fs)

	browser, _ := show.NewBrowser(ctx, "web://")

	show_opts := &show.RunOptions{
		Browser: browser,
		Mux:     mux,
	}

	err = show.RunWithOptions(ctx, show_opts)

	if err != nil {
		return fmt.Errorf("Failed to show results, %w", err)
	}

	return nil
}
