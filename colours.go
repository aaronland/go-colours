package colours

import (
	"fmt"
	"image"
)

type Colour interface {
	Name() string
	Hex() string
	// Reference() string
	// Closest() []Colour
	String() string
}

type Palette interface {
	Reference() string
	Colours() []Colour
}

type Extruder interface {
	Colours(image.Image, int) ([]Colour, error)
}

type Grid interface {
	Closest(Colour, Palette) (Colour, error)
}

type CommonColour struct {
	Colour          `json:",omitempty"`
	CommonName      string `json:"name"`
	CommonHex       string `json:"hex"`
	CommonReference string `json:"reference,omitempty"`
}

func (hc *CommonColour) Name() string {
	return hc.CommonName
}

func (hc *CommonColour) Hex() string {
	return hc.CommonHex
}

func (hc *CommonColour) Reference() string {
	return hc.CommonReference
}

func (hc *CommonColour) String() string {

	name := hc.Name()
	hex := hc.Hex()

	if name == hex {
		return hex
	}

	return fmt.Sprintf("%s (%s)", hex, name)
}

func NewHexColour(hex_value string) (Colour, error) {

	hc := CommonColour{
		CommonName:      hex_value,
		CommonHex:       hex_value,
		CommonReference: "unknown",
	}

	return &hc, nil
}

type CommonPalette struct {
	Palette         `json:",omitempty"`
	CommonReference string          `json:"reference"`
	CommonColours   []*CommonColour `json:"colours,omitempty"`
}

func (p *CommonPalette) Reference() string {
	return p.CommonReference
}

func (p *CommonPalette) Colours() []Colour {

	// Y DO I NEED TO DOOOOOOOOOOOOOOOOOOO THIS???
	// Y U SO WEIRD GOOOOOOOOOOOOOOOO????
	// (20180605/thisisaaronland)

	c := make([]Colour, 0)

	for _, pc := range p.CommonColours {
		c = append(c, pc)
	}

	return c
}
