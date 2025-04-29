package palette

import (
	"context"
	"testing"
)

func TestCrayolaPalette(t *testing.T) {

	ctx := context.Background()

	plt, err := NewPalette(ctx, "crayola://")

	if err != nil {
		t.Fatalf("Failed to crayola palette, %v", err)
	}

	expected_ref := "crayola"

	if plt.Reference() != expected_ref {
		t.Fatalf("Invalid reference. Expected '%s', got '%s'", expected_ref, plt.Reference())
	}
}
