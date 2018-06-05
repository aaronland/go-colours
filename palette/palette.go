package palette

import (
	"encoding/json"
	"errors"
	"github.com/aaronland/go-colours"
	_ "log"
	"strings"
)

type CommonPaletteColour struct {
	colours.Colour
	CPName string `json:"name"`
	CPHex  string `json:"hex"`
}

func (c *CommonPaletteColour) Name() string {
	return c.CPName
}

func (c *CommonPaletteColour) Hex() string {
	return c.CPHex
}

type CommonPalette struct {
	colours.Palette `json:",omitempty"`
	PReference      string                 `json:"reference"`
	PColours        []*CommonPaletteColour `json:"colours,omitempty"`
}

func (p *CommonPalette) Reference() string {
	return p.PReference
}

func (p *CommonPalette) Colours() []colours.Colour {

	// Y DO I NEED TO DOOOOOOOOOOOOOOOOOOO THIS???
	// Y U SO WEIRD GOOOOOOOOOOOOOOOO????
	// (20180605/thisisaaronland)

	c := make([]colours.Colour, 0)

	for _, pc := range p.PColours {
		c = append(c, pc)
	}

	return c
}

func NewNamedPalette(name string, args ...interface{}) (colours.Palette, error) {

	var p colours.Palette
	var err error

	switch strings.ToUpper(name) {
	case "CSS4":
		p, err = NewCommonPalette(CSS4, args)
	default:
		err = errors.New("Invalid or unknown palette")
	}

	return p, err
}

func NewCommonPalette(data []byte, args ...interface{}) (colours.Palette, error) {

	var p CommonPalette

	err := json.Unmarshal(data, &p)

	if err != nil {
		return nil, err
	}

	return &p, nil
}
