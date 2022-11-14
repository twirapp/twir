package moderation

import "testing"

func TestIsToMuchSymbols(t *testing.T) {
	res := IsToMuchSymbols("♣♦•◘♠ qwerty", 30)
	if res {
		t.Errorf(`IsToMuchSymbols("♣♦•◘♠ qwerty", 30) = true; want false`)
	}

	res = IsToMuchSymbols("♣♦•◘♠ qwerty", 51)
	if !res {
		t.Errorf(`IsToMuchSymbols("♣♦•◘♠ qwerty", 50) = false; want true`)
	}
}
