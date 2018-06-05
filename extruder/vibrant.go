package extruder

import (
	"fmt"
	"github.com/RobCherry/vibrant"
	"github.com/aaronland/go-colours"
	"github.com/pwaller/go-hexcolor"
	"golang.org/x/image/draw"
	"image"
	_ "sort"
)

type VibrantExtruder struct {
	colours.Extruder
	max_colors uint32
}

func NewVibrantExtruder(args ...interface{}) (colours.Extruder, error) {

	v := VibrantExtruder{
		max_colors: 24,
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

	colours := make([]colours.Colour, 0)

	for _, sw := range swatches {

		if sw == nil {
			continue
		}

		rgba := sw.RGBAInt()
		r, g, b, a := rgba.RGBA()

		hex := hexcolor.RGBAToHex(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))
		str_hex := fmt.Sprintf("%s", hex)

		c, _ := colours.NewColourFromHex(str_hex)

		colours = append(colours, c)

		if limit > 0 && len(colours) == limit {
			break
		}
	}

	return colours, nil
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
