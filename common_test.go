package colours

import (
	"context"
	"testing"
)

func TestNewCommonColour(t *testing.T) {

	tests := []string{
		"common://?hex=dec453&name=dec453&ref=vibrant",
	}

	ctx := context.Background()

	for _, uri := range tests {

		_, err := NewCommonColour(ctx, uri)

		if err != nil {
			t.Fatalf("Failed to create colour for '%s', %v", uri, err)
		}
	}
}
