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

	"github.com/aaronland/go-colours/extruder"
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

	for _, path := range flag.Args() {

		f, err := os.Open(path)

		if err != nil {
			log.Fatal(err)
		}

		im, _, err := image.Decode(f)

		if err != nil {
			log.Fatal(err)
		}

		c, err := ex.Colours(im, 5)

		if err != nil {
			log.Fatal(err)
		}

		for _, c := range c {
			log.Println(c)
		}
	}
}
