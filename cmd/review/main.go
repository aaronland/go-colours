// Command-line tool to generate an HTML page (and associated assets) to review the colour extraction
// for an image using one or more extruders and one or more palettes. The application will spawn a short-lived
// web server to serve the HTML review on a random port number and open its URI in the default browser.
package main

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
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-colours"
	"github.com/aaronland/go-colours/extruder"
	"github.com/aaronland/go-colours/grid"
	"github.com/aaronland/go-colours/palette"
	"github.com/sfomuseum/go-flags/multi"
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

func main() {

	var extruder_uris multi.MultiString
	var palette_uris multi.MultiString

	var root string

	flag.Var(&extruder_uris, "extruder-uri", "Zero or more aaronland/go-colours/extruder.Extruder URIs. Default is to use simple:// and vibrant://")
	flag.Var(&palette_uris, "palette-uri", "Zero or more aaronland/go-colours/palette.Palette URIs. Default is to use css3://, css4:// and crayola://")
	flag.StringVar(&root, "root", "", "The path to a directory where images and HTML files associated with the review should be stored. If empty a new temporary directory will be created (and deleted when the application exits).")

	flag.Parse()

	ctx := context.Background()

	var abs_root string

	if root != "" {

		root_dir, err := filepath.Abs(root)

		if err != nil {
			log.Fatalf("Failed to derive root, %v", err)
		}

		abs_root = root_dir

	} else {

		root_dir, err := os.MkdirTemp("", "colours")

		if err != nil {
			log.Fatalf("Failed to create temp dir, %v", err)
		}

		defer os.RemoveAll(root_dir)
		abs_root = root_dir
	}

	if len(extruder_uris) == 0 {
		extruder_uris.Set("simple://")
		extruder_uris.Set("vibrant://")
	}

	if len(palette_uris) == 0 {
		palette_uris.Set("css3://")
		palette_uris.Set("css4://")
		palette_uris.Set("crayola://")
	}

	extruders := make([]extruder.Extruder, len(extruder_uris))

	for idx, ex_uri := range extruder_uris {

		ex, err := extruder.NewExtruder(ctx, ex_uri)

		if err != nil {
			log.Fatalf("Failed to create new '%s' extruder, %v", ex_uri, err)
		}

		extruders[idx] = ex
	}

	gr, err := grid.NewGrid(ctx, "euclidian://")

	if err != nil {
		log.Fatalf("Failed to create new grid, %v", err)
	}

	palettes := make([]palette.Palette, len(palette_uris))

	for idx, p_uri := range palette_uris {

		p, err := palette.NewPalette(ctx, p_uri)

		if err != nil {
			log.Fatalf("Failed to create '%s' palette, %v", p_uri, err)
		}

		palettes[idx] = p
	}

	index_t, err := template.New("index").Parse(index_html)

	if err != nil {
		log.Fatalf("Failed to parse template, %v", err)
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

	for _, path := range flag.Args() {

		if strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "http://") {

			fname := filepath.Base(path)

			rsp, err := http.Get(path)

			if err != nil {
				log.Fatalf("Failed to fetch %s, %v", path, err)
			}

			defer rsp.Body.Close()

			new_path := filepath.Join(abs_root, fname)
			new_wr, err := os.OpenFile(new_path, os.O_RDWR|os.O_CREATE, 0644)

			if err != nil {
				log.Fatalf("Failed to open %s for writing, %v", new_path, err)
			}

			_, err = io.Copy(new_wr, rsp.Body)

			if err != nil {
				log.Fatalf("Failed to copy %s to %s, %v", path, new_path, err)
			}

			err = new_wr.Close()

			if err != nil {
				log.Fatalf("Failed to close %s after writing, %v", new_path, err)
			}

			path = new_path
		}

		fname := filepath.Base(path)
		ext := filepath.Ext(fname)
		fname = strings.Replace(fname, ext, ".png", 1)

		r, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s for reading, %v", path, err)
		}

		defer r.Close()

		im, _, err := image.Decode(r)

		if err != nil {
			log.Fatalf("Failed to decode %s, %v", path, err)
		}

		extrusions, err := derive_colours(im)

		if err != nil {
			log.Fatalf("Failed to derive colours, %v", err)
		}

		im_c := &Image{
			URI:        fname,
			Extrusions: extrusions,
		}

		images = append(images, im_c)

		im_path := filepath.Join(abs_root, fname)

		im_wr, err := os.OpenFile(im_path, os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			log.Fatalf("Failed to open %s for writing, %v", im_path, err)
		}

		err = png.Encode(im_wr, im)

		if err != nil {
			log.Fatalf("Failed to encode %s, %v", im_path, err)
		}

		err = im_wr.Close()

		if err != nil {
			log.Fatalf("Failed to close %s after writing, %v", im_path, err)
		}

	}

	//

	index_path := filepath.Join(abs_root, "index.html")

	index_wr, err := os.OpenFile(index_path, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		log.Fatalf("Failed to open %s for writing, %v", index_path, err)
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
		log.Fatalf("Failed to encode %s, %v", index_path, err)
	}

	err = index_wr.Close()

	if err != nil {
		log.Fatalf("Failed to close %s after writing, %v", index_path, err)
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
		log.Fatalf("Failed to show results, %v", err)
	}
}
