package colours

import (
       "image"
)

type Colour interface {
     Name() string
     Hex() string
     Reference() string
     // Closest() []Colour
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

type HexColour struct {
     Colour
     HexName string
     HexColour string
     HexReference string
}

func (hc *HexColour) Name() string {
     return hc.HexName
}

func (hc *HexColour) Hex() string {
     return hc.HexColour
}

func (hc *HexColour) Reference() string {
     return hc.HexReference
}

func NewHexColour(hex_value string) (Colour, error){

     hc := HexColour{
     	HexName: hex_value,
	HexColour: hex_value,
	HexReference: "unknown",
     }

     return &hc, nil
}