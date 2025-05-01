package extruder

import (
	"context"
	"fmt"
	"image"
	"net/url"

	// "github.com/nfnt/resize"
	"github.com/aaronland/go-colours"
	"github.com/mccutchen/palettor"
)

const PALETTOR string = "palettor"

type PalettorExtruder struct {
	Extruder
}

func init() {
	ctx := context.Background()
	err := RegisterExtruder(ctx, PALETTOR, NewPalettorExtruder)
	if err != nil {
		panic(err)
	}
}

func NewPalettorColour(ctx context.Context, str_hex string) (colours.Colour, error) {

	u := url.URL{}
	u.Scheme = "common"

	q := url.Values{}
	q.Set("hex", str_hex)
	q.Set("name", PALETTOR)
	q.Set("ref", str_hex)

	u.RawQuery = q.Encode()

	return colours.NewColour(ctx, u.String())
}

func NewPalettorExtruder(ctx context.Context, uri string) (Extruder, error) {

	ex := PalettorExtruder{}
	return &ex, nil
}

func (ex *PalettorExtruder) Name() string {
	return PALETTOR
}

func (ex *PalettorExtruder) Colours(ctx context.Context, im image.Image, limit int) ([]colours.Colour, error) {

	// im = resize.Thumbnail(200, 200, im, resize.Lanczos3)
	
	max_iters := 100
	
	palette, err := palettor.Extract(limit, max_iters, im)

	if err != nil {
		return nil, fmt.Errorf("Failed to extract palette, %w", err)
	}

	results := make([]colours.Colour, limit)

	for i, c := range palette.Colors(){

		hex_value := toHexColor(c)
		colour, err := NewPalettorColour(ctx, hex_value)

		if err != nil {
			return nil, err
		}

		results[i] = colour
	}

	return results, nil
}
