package moderation_helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTooMuchCaps(t *testing.T) {
	table := []struct {
		msg           string
		want          bool
		maxPercentage int
	}{
		{
			msg:           "HELLO",
			want:          true,
			maxPercentage: 50,
		},
		{
			msg:           "HELLO",
			want:          false,
			maxPercentage: 100,
		},
		{
			msg:           "QWERTYuiop",
			want:          true,
			maxPercentage: 50,
		},
		{
			msg:           "QWERTyuiop",
			want:          false,
			maxPercentage: 51,
		},
	}

	for _, tt := range table {
		t.Run(
			tt.msg,
			func(t *testing.T) {
				ok, capsCount := IsTooMuchCaps(tt.msg, tt.maxPercentage)
				assert.Equal(
					t,
					tt.want,
					ok,
					"IsTooMuchCaps(%q, %d) = %v, want %v, capsLength %v",
					tt.msg,
					tt.maxPercentage,
					ok,
					tt.want,
					capsCount,
				)
			},
		)
	}
}
