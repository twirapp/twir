package moderation

import "testing"

func TestIsTooLong(t *testing.T) {
	cases := []struct {
		name      string
		msg       string
		maxLength int
		expected  bool
	}{
		{name: "v", msg: "v", maxLength: 1, expected: false},
		{name: "vk", msg: "vk", maxLength: 1, expected: true},
		{name: "vk", msg: "vk", maxLength: 3, expected: false},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got := IsTooLong(test.msg, test.maxLength)
			if got != test.expected {
				t.Errorf(
					"msg=%q maxLength=%v expected=%v but got=%v",
					test.msg,
					test.maxLength,
					test.expected,
					got,
				)
			}
		})
	}
}
