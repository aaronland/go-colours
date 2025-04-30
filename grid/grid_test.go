package grid

import (
	"context"
	"testing"
)

func TestNewGrid(t *testing.T) {

	tests := []string{
		"euclidian://",
	}

	ctx := context.Background()

	for _, uri := range tests {

		_, err := NewGrid(ctx, uri)

		if err != nil {
			t.Fatalf("Failed to create grid for %s, %v", uri, err)
		}
	}
}
