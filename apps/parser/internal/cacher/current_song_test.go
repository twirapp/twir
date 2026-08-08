package cacher

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsSpotifyInvalidGrantError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "expired refresh token",
			err: errors.New(
				`decode spotify refresh response: status 400: {"error":"invalid_grant","error_description":"Refresh token expired"}`,
			),
			want: true,
		},
		{
			name: "revoked refresh token",
			err: errors.New(
				`decode spotify refresh response: status 400: {"error":"invalid_grant","error_description":"Refresh token revoked"}`,
			),
			want: true,
		},
		{
			name: "spotify server error",
			err:  errors.New(`decode spotify refresh response: status 500: {"error":"server_error"}`),
			want: false,
		},
		{
			name: "network error",
			err:  errors.New("refresh spotify token: connection refused"),
			want: false,
		},
		{
			name: "bus timeout",
			err:  errors.New("timeout"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, isSpotifyInvalidGrantError(tt.err))
		})
	}
}
