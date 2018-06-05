package palette

import (
	"encoding/json"
	"errors"
	"github.com/aaronland/go-colours"
	"strings"
)

type CommonPalette struct {
	colours.Palette
	PReference string           `json:"reference"`
	PColours   []colours.Colour `json:colours"`
}

func (p *CommonPalette) Reference() string {
	return p.PReference
}

func (p *CommonPalette) Colours() []colours.Colour {
	return p.PColours
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
