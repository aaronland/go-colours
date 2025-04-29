package extruder

import (
	"context"
	"fmt"
	"image"
	"strings"

	"github.com/RobCherry/vibrant"
	"github.com/aaronland/go-colours"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/draw"
)

type VibrantExtruder struct {
	Extruder
	max_colours uint32
}

func init() {
	ctx := context.Background()
	err := RegisterExtruder(ctx, "vibrant", NewVibrantExtruder)
	if err != nil {
		panic(err)
	}
}

func NewVibrantExtruder(ctx context.Context, uri string) (Extruder, error) {

	v := VibrantExtruder{
		max_colours: 24,
	}

	return &v, nil
}

func (v *VibrantExtruder) Colours(im image.Image, limit int) ([]colours.Colour, error) {

	pb := vibrant.NewPaletteBuilder(im)
	pb = pb.MaximumColorCount(v.max_colours)
	pb = pb.Scaler(draw.ApproxBiLinear)

	palette := pb.Generate()

	// swatches := palette.Swatches()
	// sort.Sort(populationSwatchSorter(swatches))

	swatches := []*vibrant.Swatch{
		palette.VibrantSwatch(),
		palette.LightVibrantSwatch(),
		palette.DarkVibrantSwatch(),
		palette.MutedSwatch(),
		palette.LightMutedSwatch(),
		palette.DarkMutedSwatch(),
	}

	results := make([]colours.Colour, 0)

	for _, sw := range swatches {

		if sw == nil {
			continue
		}

		cl, ok := colorful.MakeColor(sw.Color())

		if !ok {
			return nil, fmt.Errorf("Unable to make color, %v", sw.Color())
		}

		hex := cl.Hex()
		hex = strings.TrimLeft(hex, "#")

		ctx := context.Background()

		c_uri := fmt.Sprintf("common://?hex=%s&name=%s&ref=vibrant", hex, hex)
		c, err := colours.NewColour(ctx, c_uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to create new color '%s', %w", c_uri, err)
		}

		results = append(results, c)

		if limit > 0 && len(results) == limit {
			break
		}
	}

	return results, nil
}

// these are straight copies of vibrant/cli/main.go

type populationSwatchSorter []*vibrant.Swatch

func (p populationSwatchSorter) Len() int           { return len(p) }
func (p populationSwatchSorter) Less(i, j int) bool { return p[i].Population() > p[j].Population() }
func (p populationSwatchSorter) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type hueSwatchSorter []*vibrant.Swatch

func (p hueSwatchSorter) Len() int           { return len(p) }
func (p hueSwatchSorter) Less(i, j int) bool { return p[i].HSL().H < p[j].HSL().H }
func (p hueSwatchSorter) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
