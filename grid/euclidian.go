package grid

import (
	"errors"
	"github.com/aaronland/go-colours"
	"github.com/pwaller/go-hexcolor"
	"math"
	"sort"
)

type EuclidianGrid struct {
	colours.Grid
}

func NewEuclidianGrid(args ...interface{}) (colours.Grid, error) {

	eu := EuclidianGrid{}

	return &eu, nil
}

func (eu *EuclidianGrid) Closest(target colours.Colour, palette colours.Palette) (colours.Colour, error) {

	// https://github.com/pwaller/go-hexcolor/blob/master/hexcolor.go
	// https://github.com/ubernostrum/webcolors/blob/master/webcolors.py#L473-L485

	r, g, b, _ := hexcolor.HexToRGBA(hexcolor.Hex(target.Hex()))

	lookup := make(map[int]colours.Colour)
	keys := make([]int, 0)

	for _, c := range palette.Colours() {

		rc, gc, bc, _ := hexcolor.HexToRGBA(hexcolor.Hex(c.Hex()))

		rd := math.Pow(float64(int32(rc)-int32(r)), 2.0)
		gd := math.Pow(float64(int32(gc)-int32(g)), 2.0)
		bd := math.Pow(float64(int32(bc)-int32(b)), 2.0)

		k := int(rd + gd + bd)
		lookup[k] = c

		keys = append(keys, k)
	}

	sort.Ints(keys)

	/*
		for i, idx := range keys {
			log.Println(i, idx, lookup[idx])
		}
	*/

	if len(keys) == 0 {
		return nil, errors.New("Nothing found")
	}

	return lookup[keys[0]], nil
}
