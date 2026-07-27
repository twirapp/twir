package quotes

import "testing"

func TestParseQuoteNumber_whenInputIsValid(t *testing.T) {
	cases := []struct {
		input    string
		expected int
	}{
		{input: "42", expected: 42},
		{input: "#42", expected: 42},
		{input: " #42 ", expected: 42},
	}

	for _, testCase := range cases {
		t.Run(testCase.input, func(t *testing.T) {
			number, ok := parseQuoteNumber(testCase.input)
			if !ok {
				t.Fatal("expected quote number to parse")
			}
			if number != testCase.expected {
				t.Fatalf("expected %d, got %d", testCase.expected, number)
			}
		})
	}
}

func TestParseQuoteNumber_whenInputIsInvalid(t *testing.T) {
	for _, input := range []string{"", "#", "0", "#0", "nope"} {
		t.Run(input, func(t *testing.T) {
			_, ok := parseQuoteNumber(input)
			if ok {
				t.Fatal("expected quote number to be invalid")
			}
		})
	}
}
