// Command-line tool to generate an HTML page (and associated assets) to review the colour extraction
// for an image using one or more extruders and one or more palettes. The application will spawn a short-lived
// web server to serve the HTML review on a random port number and open its URI in the default browser.
package main

import (
	"context"
	"log"

	"github.com/aaronland/go-colours/app/review"
)

func main() {

	ctx := context.Background()
	err := review.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
