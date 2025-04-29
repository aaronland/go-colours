package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/aaronland/go-colours"
	"github.com/aaronland/go-colours/grid"
	"github.com/aaronland/go-colours/palette"
)

func main() {

	var grid_uri string
	var palette_uri string

	flag.StringVar(&grid_uri, "grid-uri", "euclidian://", "...")
	flag.StringVar(&palette_uri, "palette-uri", "css3://", "...")

	flag.Parse()

	ctx := context.Background()

	gr, err := grid.NewGrid(ctx, grid_uri)

	if err != nil {
		log.Fatal(err)
	}

	p, err := palette.NewPalette(ctx, palette_uri)

	if err != nil {
		log.Fatal(err)
	}

	for _, hex := range flag.Args() {

		c_uri := fmt.Sprintf("common://?hex=%s", hex)

		target, err := colours.NewColour(ctx, c_uri)

		if err != nil {
			log.Fatal(err)
		}

		match, err := gr.Closest(target, p)

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%s SNAPS TO %s\n", target, match)
	}
}
