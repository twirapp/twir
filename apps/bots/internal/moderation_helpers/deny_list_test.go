package moderation_helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasDeniedWord(t *testing.T) {
	table := []struct {
		msg  string
		list []string
		want bool
	}{
		{
			msg:  "hello",
			list: []string{"hello"},
			want: true,
		},
		{
			msg:  "hello",
			list: []string{"hello", "world"},
			want: true,
		},
		{
			msg:  "helloworld",
			list: []string{"(?m)hello"},
			want: true,
		},
		{
			msg:  "hell",
			list: []string{"(?m)hello"},
			want: false,
		},
		{
			msg:  "hello world",
			list: []string{"hello", "world"},
			want: true,
		},
		{
			msg:  "hello world",
			list: []string{},
			want: false,
		},
		{
			msg:  "hello world",
			list: []string{"foo"},
			want: false,
		},
	}

	for _, tt := range table {
		t.Run(
			tt.msg, func(t *testing.T) {
				got := HasDeniedWord(tt.msg, tt.list)
				assert.Equal(
					t,
					tt.want,
					got,
					"HasDeniedWord(%q, %q) = %v, want %v",
					tt.msg,
					tt.list,
					got,
					tt.want,
				)
			},
		)

	}
}
