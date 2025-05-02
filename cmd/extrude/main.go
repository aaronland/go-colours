// Command line tool to extrude (derive) dominant colours from one or more images as well as closest matches colours
// using zero or more "snap-to-grid" colour palettes as JSON-encoded data written to STDOUT.
package main

import (
	"context"
	"log"

	"github.com/aaronland/go-colours/app/extrude"
)

func main() {

	ctx := context.Background()
	err := extrude.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
