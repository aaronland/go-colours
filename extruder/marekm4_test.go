package extruder

import (
	"context"
	"testing"
)

func TestSimpleExtruder(t *testing.T) {

	ctx := context.Background()

	_, err := NewExtruder(ctx, "marekm4://")

	if err != nil {
		t.Fatalf("Failed to marekm4 extruder, %v", err)
	}

	// Do stuff here...
}
