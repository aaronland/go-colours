package extruder

import (
	"context"
	"testing"
)

func TestSimpleExtruder(t *testing.T) {

	ctx := context.Background()

	_, err := NewExtruder(ctx, "simple://")

	if err != nil {
		t.Fatalf("Failed to vibrant extruder, %v", err)
	}

	// Do stuff here...
}
