package palette

import (
	"context"
	"testing"
)

func TestNewPalette(t *testing.T) {

	tests := []string{
		"crayola://",
		"css3://",
		"css4://",
	}

	ctx := context.Background()

	for _, uri := range tests {

		_, err := NewPalette(ctx, uri)

		if err != nil {
			t.Fatalf("Failed to create palette for %s, %v", uri, err)
		}
	}
}
