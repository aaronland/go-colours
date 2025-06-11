package extruder

import (
	"context"
	"testing"
)

func TestNewMarekm4Colour(t *testing.T) {

	ctx := context.Background()
	hex := "#cc6699"

	c, err := NewMarekm4Colour(ctx, hex)

	if err != nil {
		t.Fatalf("Failed to create new marekm4 colour, %v", err)
	}

	if c.Reference() != MAREKM4 {
		t.Fatalf("Invalid ref for marekm4 colour, %s", c.Reference())
	}

	if c.Hex() != hex {
		t.Fatalf("Invalid hex, %s (expected %s)", c.Hex(), hex)
	}

	if c.Name() != hex {
		t.Fatalf("Invalid name, %s (expected %s)", c.Name(), hex)
	}

}

func TestMarekm4Extruder(t *testing.T) {

	ctx := context.Background()

	_, err := NewExtruder(ctx, "marekm4://")

	if err != nil {
		t.Fatalf("Failed to marekm4 extruder, %v", err)
	}

	// Do stuff here...
}
