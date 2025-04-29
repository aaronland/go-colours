package grid

import (
	"context"
	"testing"
)

func TestSimpleGrid(t *testing.T) {

	ctx := context.Background()

	_, err := NewGrid(ctx, "euclidian://")

	if err != nil {
		t.Fatalf("Failed to euclidian grid, %v", err)
	}

	// Do stuff here...
}
