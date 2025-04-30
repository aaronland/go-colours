package grid

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aaronland/go-colours"
	"github.com/aaronland/go-colours/palette"
)

func TestEuclidiagGrid(t *testing.T) {

	ctx := context.Background()

	_, err := NewGrid(ctx, "euclidian://")

	if err != nil {
		t.Fatalf("Failed to euclidian grid, %v", err)
	}

	// Do stuff here...
}

func TestEuclidianClosestCSS3(t *testing.T) {

	tests := map[string]string{
		"#cc3366": "#cd5c5c",
	}

	ctx := context.Background()

	gr, err := NewGrid(ctx, "euclidian://")

	if err != nil {
		t.Fatalf("Failed to euclidian grid, %v", err)
	}

	plt, err := palette.NewPalette(ctx, "css3://")

	if err != nil {
		t.Fatalf("Failed to create css3 palette, %v", err)
	}

	for hex, expected := range tests {

		hex = strings.TrimLeft(hex, "#")

		c_uri := fmt.Sprintf("common://?hex=%s", hex)

		target, err := colours.NewColour(ctx, c_uri)

		if err != nil {
			t.Fatalf("Failed to create new colour for uri '%s', %v", c_uri, err)
		}

		match, err := gr.Closest(target, plt)

		if err != nil {
			t.Fatalf("Failed to derive closest match, %v", err)
		}

		match_hex := match.Hex()

		if match_hex != expected {
			t.Fatalf("Invalid match for '%s'. Expected '%s' but got '%s'", hex, expected, match_hex)
		}
	}

}
