package palette

import (
	"context"
	"testing"
)

func TestCSS3Palette(t *testing.T) {

	ctx := context.Background()

	plt, err := NewPalette(ctx, "css3://")

	if err != nil {
		t.Fatalf("Failed to css3 palette, %v", err)
	}

	expected_ref := "css3"

	if plt.Reference() != expected_ref {
		t.Fatalf("Invalid reference. Expected '%s', got '%s'", expected_ref, plt.Reference())
	}
}
