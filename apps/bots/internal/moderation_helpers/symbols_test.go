package moderation_helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsToMuchSymbols(t *testing.T) {
	type args struct {
		msg           string
		maxPercentage int
	}
	tests := []struct {
		name      string
		args      args
		want      bool
		wantCount int
	}{
		{
			name: "should be false", args: args{msg: "test", maxPercentage: 10}, want: false,
			wantCount: 0,
		},
		{name: "should be true", args: args{msg: "..", maxPercentage: 1}, want: true, wantCount: 2},
		{name: "should be true", args: args{msg: ".test.", maxPercentage: 1}, want: true, wantCount: 2},
		{name: "should be true", args: args{msg: "⣿", maxPercentage: 1}, want: true, wantCount: 1},
		{name: "should be true", args: args{msg: "👉🏿👈🏿 ", maxPercentage: 1}, want: true, wantCount: 4},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ok, count := IsToMuchSymbols(tt.args.msg, tt.args.maxPercentage)

				assert.Equal(
					t,
					tt.want,
					ok,
					"IsToMuchSymbols(%v, %v)",
					tt.args.msg,
					tt.args.maxPercentage,
				)
				assert.Equal(
					t,
					tt.wantCount,
					count,
					"IsToMuchSymbols(%v, %v)",
					tt.args.msg,
					tt.args.maxPercentage,
				)
			},
		)
	}
}
