package main

import (
	"flag"
	"github.com/aaronland/go-colours/extract"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

func main() {

	flag.Parse()

	for _, path := range flag.Args() {

		f, err := os.Open(path)

		if err != nil {
			log.Fatal(err)
		}

		im, _, err := image.Decode(f)

		if err != nil {
			log.Fatal(err)
		}

		c, err := extract.Extract(im)

		log.Println(c)
	}
}
