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

	flag.StringVar(&extruder_uri, "extruder-uri", "vibrant://", "...")

	flag.Parse()

	ctx := context.Background()

	ex, err := extruder.NewExtruder(ctx, extruder_uri)

	if err != nil {
		log.Fatalf("Failed to create new extruder, %v", err)
	}

	for _, path := range flag.Args() {

		f, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s for reading, %v", path, err)
		}

		im, _, err := image.Decode(f)

		if err != nil {
			log.Fatalf("Failed to decode %s, %v", path, err)
		}

		c, err := ex.Colours(im, 5)

		if err != nil {
			log.Fatalf("Failed to derive colours, %v", err)
		}

		for _, c := range c {
			log.Println(c)
		}
	}
}
