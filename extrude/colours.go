package extrude

import (
	"fmt"
	"sync"

	"github.com/aaronland/go-colours"
)

func UniqueColours(images []*Image) []colours.Colour {

	lookup := new(sync.Map)
	derived := make([]colours.Colour, 0)

	add_colour := func(c colours.Colour) {

		k := fmt.Sprintf("%s#%s", c.Reference(), c.Hex())

		_, exists := lookup.LoadOrStore(k, true)

		if !exists {
			derived = append(derived, c)
		}
	}

	for _, im := range images {

		for _, ex := range im.Extrusions {

			for _, sw := range ex.Swatches {

				add_colour(sw.Colour)

				for _, cl := range sw.Closest {
					add_colour(cl.Colour)
				}
			}
		}
	}

	return derived
}
