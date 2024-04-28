package command_arguments

import (
	"testing"
)

func TestParser(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		args          []Arg
		input         string
		expectErr     bool
		checkValidity func(*Parser) error
	}{
		{
			name:      "empty args",
			args:      []Arg{},
			input:     "",
			expectErr: false,
		},
		{
			name:      "empty input",
			args:      []Arg{String{Name: "test"}},
			input:     "",
			expectErr: true,
		},
		{
			name:      "invalid string arg",
			args:      []Arg{Int{Name: "test"}},
			input:     "test",
			expectErr: true,
		},
		{
			name:      "invalid int arg",
			args:      []Arg{Int{Name: "test"}},
			input:     "test",
			expectErr: true,
		},
		{
			name:      "valid arg",
			args:      []Arg{Int{Name: "test"}},
			input:     "1",
			expectErr: false,
		},
		{
			name:      "optional arg",
			args:      []Arg{Int{Name: "test", Optional: true}},
			input:     "",
			expectErr: false,
		},
		{
			name:      "optional arg with value",
			args:      []Arg{Int{Name: "test", Optional: true}},
			input:     "1",
			expectErr: false,
		},
		{
			name:      "variadic string",
			args:      []Arg{VariadicString{Name: "test"}},
			input:     "test",
			expectErr: false,
			checkValidity: func(p *Parser) error {
				if p.Get("test").String() != "test" {
					return ErrInvalidArg
				}

				return nil
			},
		},
		{
			name:      "variadic string with multiple words",
			args:      []Arg{VariadicString{Name: "test"}},
			input:     "test test",
			expectErr: false,
			checkValidity: func(p *Parser) error {
				if p.Get("test").String() != "test test" {
					return ErrInvalidArg
				}

				return nil
			},
		},
		{
			name:      "int",
			args:      []Arg{Int{Name: "test"}},
			input:     "1",
			expectErr: false,
			checkValidity: func(p *Parser) error {
				if p.Get("test").Int() != 1 {
					return ErrInvalidArg
				}

				return nil
			},
		},
		{
			name:      "optional should be nil",
			args:      []Arg{Int{Name: "test", Optional: true}},
			input:     "",
			expectErr: false,
			checkValidity: func(p *Parser) error {
				if p.Get("test") != nil {
					return ErrInvalidArg
				}

				return nil
			},
		},
	}

	for _, c := range cases {
		t.Run(
			c.name, func(t *testing.T) {
				parser, err := NewParser(c.args, c.input)
				if err != nil && !c.expectErr {
					t.Errorf("expected error: %v, got: %v", c.expectErr, err)
				}

				if c.checkValidity != nil {
					err := c.checkValidity(parser)
					if err != nil {
						t.Errorf("unexpected error: %v", err)
					}
				}
			},
		)
	}
}
