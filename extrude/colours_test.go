package extrude

import (
	"context"
	"testing"
)

func TestUniqueColours(t *testing.T) {

	ctx := context.Background()

	uri := "../fixtures/1762704051_KParoWjtF6WapobIDM40yA6hXOlHWTvX_b.jpg"

	opts := &ExtrudeOptions{
		Images: []string{
			uri,
		},
		ExtruderURIs: []string{
			"marekm4://",
			"vibrant://",
		},
		PaletteURIs: []string{
			"css4://",
			"crayola://",
		},
	}

	rsp, err := Extrude(ctx, opts)

	if err != nil {
		t.Fatalf("Failed to extrude colours for image, %v", err)
	}

	c := UniqueColours(rsp.Images)

	if len(c) == 0 {
		t.Fatalf("Failed to derive unique colours for image")
	}
}
