package extruder

import (
	"context"
	"testing"
	"image/color"
)

func TestVibrantExtruder(t *testing.T) {

	ctx := context.Background()

	_, err := NewExtruder(ctx, "vibrant://")

	if err != nil {
		t.Fatalf("Failed to vibrant extruder, %v", err)
	}

	// Do stuff here...
}

func TestIsTransparentFilter(t *testing.T) {

	f := new(IsTransparentFilter)
	
	c := color.NRGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 0,
	}

	if f.IsAllowed(c){
		t.Fatalf("colour should NOT be allowed")
	}

	c2 := color.NRGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 10,
	}

	if !f.IsAllowed(c2){
		t.Fatalf("colour SHOULD be allowed")
	}
	
}
