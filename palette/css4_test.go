package palette

import (
	"context"
	"testing"
)

func TestCSS4Palette(t *testing.T) {

	ctx := context.Background()

	plt, err := NewPalette(ctx, "css4://")

	if err != nil {
		t.Fatalf("Failed to css4 palette, %v", err)
	}

	expected_ref := "css4"

	if plt.Reference() != expected_ref {
		t.Fatalf("Invalid reference. Expected '%s', got '%s'", expected_ref, plt.Reference())
	}
}
