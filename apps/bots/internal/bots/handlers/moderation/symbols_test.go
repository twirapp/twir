package moderation

import "testing"

func TestIsToMuchSymbols(t *testing.T) {
	for _, test := range []struct {
		name          string
		msg           string
		maxPercentage int
		expected      bool
	}{
		{name: "false case", msg: "♣♦•◘♠ qwerty", maxPercentage: 30, expected: false},
		{name: "true case", msg: "♣♦•◘♠ qwerty", maxPercentage: 51, expected: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := IsToMuchSymbols(test.msg, test.maxPercentage)
			if got != test.expected {
				t.Errorf("msg=%q expected=%v but got=%v", test.msg, test.expected, got)
			}
		})
	}
}
