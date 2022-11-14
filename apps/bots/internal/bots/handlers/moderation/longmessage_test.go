package moderation

import "testing"

func TestIsTooLong(t *testing.T) {
	if !IsTooLong("qwerty", 5) {
		t.Errorf(`IsTooLong("qwerty", 5) = %v; want true`, false)
	}
}
