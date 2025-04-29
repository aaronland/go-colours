package main

import (
	"context"
	"flag"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	_ "github.com/aaronland/go-colours"
	"github.com/aaronland/go-colours/extruder"
	"github.com/aaronland/go-colours/grid"
	"github.com/aaronland/go-colours/palette"
)

func main() {

	var extruder_uri string

	flag.StringVar(&extruder_uri, "extruder-uri", "virbant://", "...")

	flag.Parse()

	ctx := context.Background()

	ex, err := extruder.NewExtruder(ctx, extruder_uri)

	if err != nil {
		log.Fatal(err)
	}

	gr, err := grid.NewNamedGrid("euclidian")

	if err != nil {
		log.Fatal(err)
	}

	p, err := palette.NewNamedPalette("css4")

	if err != nil {
		log.Fatal(err)
	}

	for _, path := range flag.Args() {

		f, err := os.Open(path)

		if err != nil {
			log.Fatal(err)
		}

		im, _, err := image.Decode(f)

		if err != nil {
			log.Fatal(err)
		}

		log.Println(path)

		c, err := ex.Colours(im, 5)

		if err != nil {
			log.Fatal(err)
		}

		for _, c := range c {
			log.Println(c)

			cl, _ := gr.Closest(c, p)

			log.Println(cl)
		}

	}
}
