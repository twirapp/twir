package shoutout

import "testing"

func TestHasShoutoutScope_returns_true_only_for_moderator_manage_shoutouts(t *testing.T) {
	// Given
	tests := []struct {
		name   string
		scopes []string
		want   bool
	}{
		{name: "required scope", scopes: []string{"moderator:manage:shoutouts"}, want: true},
		{name: "other moderator scope", scopes: []string{"moderator:read:followers"}, want: false},
		{name: "empty scopes", scopes: nil, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// When
			got := hasShoutoutScope(test.scopes)

			// Then
			if got != test.want {
				t.Fatalf("hasShoutoutScope(%v) = %t, want %t", test.scopes, got, test.want)
			}
		})
	}
}
