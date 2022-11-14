package moderation

import "testing"

func TestIsTooMuchCaps(t *testing.T) {
	for _, test := range []struct {
		name          string
		msg           string
		maxPercentage int
		expected      bool
	}{
		{name: "false case", msg: "QWERTyuiop", maxPercentage: 60, expected: false},
		{name: "true case", msg: "QWERTyuiop", maxPercentage: 50, expected: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := IsTooMuchCaps(test.msg, test.maxPercentage)
			if got != test.expected {
				t.Errorf("msg=%q expected=%v but got=%v", test.msg, test.expected, got)
			}
		})
	}
}
