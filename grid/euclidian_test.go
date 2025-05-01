package grid

import (
	"context"
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

		target_c, err := colours.NewHexColour(ctx, hex)

		if err != nil {
			t.Fatalf("Failed to create new colour for hex '%s', %v", hex, err)
		}

		match, err := gr.Closest(ctx, target_c, plt)

		if err != nil {
			t.Fatalf("Failed to derive closest match, %v", err)
		}

		match_hex := match.Hex()

		if match_hex != expected {
			t.Fatalf("Invalid match for '%s'. Expected '%s' but got '%s'", hex, expected, match_hex)
		}
	}

}
