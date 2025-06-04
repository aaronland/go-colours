package extrude

import (
	"context"
	"fmt"
	"image"

	"github.com/aaronland/go-colours/extruder"
	"github.com/aaronland/go-colours/grid"
	"github.com/aaronland/go-colours/palette"
)

type DeriveExtrusionsOptions struct {
	Grid      grid.Grid
	Palettes  []palette.Palette
	Extruders []extruder.Extruder
}

func DeriveExtrusions(ctx context.Context, opts *DeriveExtrusionsOptions, im image.Image) ([]*Extrusion, error) {

	extrusions := make([]*Extrusion, 0)

	for _, ex := range opts.Extruders {

		swatches := make([]*Swatch, 0)

		colours, err := ex.Colours(ctx, im, 5)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive colours for image, %w", err)
		}

		for _, c := range colours {

			closest := make([]*Closest, 0)

			for _, p := range opts.Palettes {

				c2, err := opts.Grid.Closest(ctx, c, p)

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

		for _, p := range opts.Palettes {
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
