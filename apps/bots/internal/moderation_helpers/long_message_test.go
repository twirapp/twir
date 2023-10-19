package moderation_helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTooLong(t *testing.T) {
	tests := []struct {
		input     string
		maxLength int
		want      bool
	}{
		{input: "hello", maxLength: 1, want: true},
		{input: "empty", maxLength: 10, want: false},
	}
	for _, tt := range tests {
		t.Run(
			tt.input, func(t *testing.T) {
				assert.Equalf(
					t,
					tt.want,
					IsTooLong(tt.input, tt.maxLength),
					"IsTooLong(%v, %v)",
					tt.input,
					tt.maxLength,
				)
			},
		)
	}
}
