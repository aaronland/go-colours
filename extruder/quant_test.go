package extruder

import (
	"context"
	"testing"
)

func TestNewQuantColour(t *testing.T) {

	ctx := context.Background()
	hex := "#cc6699"

	c, err := NewQuantColour(ctx, hex)

	if err != nil {
		t.Fatalf("Failed to create new quant colour, %v", err)
	}

	if c.Reference() != QUANT {
		t.Fatalf("Invalid ref for quant colour, %s", c.Reference())
	}

	if c.Hex() != hex {
		t.Fatalf("Invalid hex, %s (expected %s)", c.Hex(), hex)
	}

	if c.Name() != hex {
		t.Fatalf("Invalid name, %s (expected %s)", c.Name(), hex)
	}

}

func TestQuantExtruder(t *testing.T) {

	ctx := context.Background()

	_, err := NewExtruder(ctx, "quant://")

	if err != nil {
		t.Fatalf("Failed to quant extruder, %v", err)
	}

	// Do stuff here...
}
