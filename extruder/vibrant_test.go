package extruder

import (
	"context"
	"testing"
)

func TestVibrantExtruder(t *testing.T) {

	ctx := context.Background()

	_, err := NewExtruder(ctx, "vibrant://")

	if err != nil {
		t.Fatalf("Failed to vibrant extruder, %v", err)
	}

	// Do stuff here...
}
