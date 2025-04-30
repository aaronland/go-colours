# go-colours

Go package for working with colours, principally colour extraction and "snap to grid"

## Important

This is work in progress. It appears to have bugs.

## Documentation

Documentation is incomplete.

## Example

```
package main

import (
	"context"
	"flag"
	"image"
	_ "image/jpeg"
	"log"
	"os"

	"github.com/aaronland/go-colours/extruder"
	"github.com/aaronland/go-colours/grid"
	"github.com/aaronland/go-colours/palette"
)

func main() {

	flag.Parse()

	ex, _ := extruder.NewExtruder(ctx, "vibrant://")
	gr, _ := grid.NewGrid(ctx, "euclidian://")
	p, _ := palette.NewPalette(ctx, "css4://")

	for _, path := range flag.Args() {

		fh, _ := os.Open(path)
		im, _, _ := image.Decode(fh)

		colours, _ := ex.Colours(im, 5)

		for _, c := range colours {

			closest, _ := gr.Closest(c, p)

			for _, cl := range closest {
				log.Println(c, cl)
			}
		}

	}
}
```

_Note that error handling has been removed for the sake of brevity._

## Interfaces

### Colour

```
type Colour interface {
	Name() string
	Hex() string
	Reference() string
	Closest() []Colour
	AppendClosest(Colour) error // I don't love this... (20180605/thisisaaronland)
	String() string
}
```

### Extruder

```
type Extruder interface {
	Colours(image.Image, int) ([]Colour, error)
}
```

### Grid

```
type Grid interface {
	Closest(Colour, Palette) (Colour, error)
}
```

### Palette

```
type Palette interface {
	Reference() string
	Colours() []Colour
}
```

## Extruders

Extruders are the things that generate a palette of colours for an `image.Image`.

### vibrant

This returns colours using the [vibrant](github.com/RobCherry/vibrant) package but rather than ranking colours using a particular metric it returns specific named "swatches" that are recast as `colours.Colour` interfaces. They are: `VibrantSwatch, LightVibrantSwatch, DarkVibrantSwatch, MutedSwatch, LightMutedSwatch, DarkMutedSwatch`.

### Grids

Grids are the things that perform operations or compare colours.

### euclidian

### Palettes

Palettes are a fixed set of colours.

### crayola

### css3

### css4

## See also

* https://github.com/RobCherry/vibrant
* https://github.com/lucasb-eyer/go-colorful

* https://github.com/aaronland/py-cooperhewitt-swatchbook
* https://github.com/aaronland/py-cooperhewitt-roboteyes-colors
* https://github.com/givp/RoyGBiv
