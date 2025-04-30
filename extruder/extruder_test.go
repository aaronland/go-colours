package extruder

import (
	"context"
	"testing"
)

func TestNewExtruder(t *testing.T) {

	tests := []string{
		"simple://",
		"vibrant://",
	}

	ctx := context.Background()

	for _, uri := range tests {

		_, err := NewExtruder(ctx, uri)

		if err != nil {
			t.Fatalf("Failed to create extruder for %s, %v", uri, err)
		}
	}
}
