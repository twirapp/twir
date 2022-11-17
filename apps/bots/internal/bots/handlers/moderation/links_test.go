package moderation

import "testing"

func TestHasLink(t *testing.T) {
	cases := []struct {
		name       string
		msg        string
		withSpaces bool
		expected   bool
	}{
		{name: "vk.com", msg: "vk.com", withSpaces: true, expected: true},
		{name: "vk. com", msg: "vk. com", withSpaces: true, expected: true},
		{name: "no link", msg: "hello world", withSpaces: true, expected: false},
		{name: "test.vscode", msg: "test.vscode", withSpaces: true, expected: false},

		{name: "vk.com", msg: "vk.com", withSpaces: false, expected: true},
		{name: "vk. com", msg: "vk. com", withSpaces: false, expected: false},
		{name: "no link", msg: "hello world", withSpaces: false, expected: false},
		{name: "test.vscode", msg: "test.vscode", withSpaces: false, expected: false},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got := HasLink(test.msg, test.withSpaces)
			if got != test.expected {
				t.Errorf(
					"msg=%q withspaces=%v expected=%v but got=%v",
					test.msg,
					test.withSpaces,
					test.expected,
					got,
				)
			}
		})
	}
}
