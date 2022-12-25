package moderation

import "testing"

func TestHasBlackListedWord(t *testing.T) {
	cases := []struct {
		name     string
		msg      string
		list     []string
		expected bool
	}{
		{name: "search vk", msg: "vk.com", list: []string{"vk"}, expected: true},
		{name: "search com", msg: "vk.com", list: []string{"com"}, expected: true},
		{name: "search duck with spaces", msg: "hi duck", list: []string{"duck"}, expected: true},
		{name: "search duck without spaces", msg: "hiduck", list: []string{"duck"}, expected: true},
		{name: "search nothing", msg: "hi duck", list: []string{}, expected: false},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got := HasBlackListedWord(test.msg, test.list)
			if got != test.expected {
				t.Errorf(
					"msg=%q withspaces=%v expected=%v but got=%v",
					test.msg,
					test.list,
					test.expected,
					got,
				)
			}
		})
	}
}
