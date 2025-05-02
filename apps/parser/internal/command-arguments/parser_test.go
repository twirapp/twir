package command_arguments

import (
	"testing"
)

func TestParser(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name                string
		args                []Arg
		input               string
		delimiter           string
		expectErr           bool
		expectedArgsResults map[string]any
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
			expectedArgsResults: map[string]any{
				"test": "test",
			},
		},
		{
			name:      "variadic string with multiple words",
			args:      []Arg{VariadicString{Name: "test"}},
			input:     "test test",
			expectErr: false,
			expectedArgsResults: map[string]any{
				"test": "test test",
			},
		},
		{
			name:      "int",
			args:      []Arg{Int{Name: "test"}},
			input:     "1",
			expectErr: false,
			expectedArgsResults: map[string]any{
				"test": 1,
			},
		},
		{
			name:      "optional should be nil",
			args:      []Arg{Int{Name: "test", Optional: true}},
			input:     "",
			expectErr: false,
			expectedArgsResults: map[string]any{
				"test": nil,
			},
		},
		{
			name: "multiple correct args",
			args: []Arg{
				Int{Name: "test1"},
				String{Name: "test2"},
			},
			input:     "1 test",
			expectErr: false,
			expectedArgsResults: map[string]any{
				"test1": 1,
				"test2": "test",
			},
		},
		{
			name: "multiple correct args with optional",
			args: []Arg{
				Int{Name: "test1"},
				String{Name: "test2", Optional: true},
			},
			input:     "1",
			expectErr: false,
		},
		{
			name: "multiple args with incorrect arg",
			args: []Arg{
				String{Name: "test1"},
				Int{Name: "test2"},
			},
			input:     "test qwe",
			expectErr: true,
		},
		{
			name: "custom delimiter",
			args: []Arg{
				String{Name: "first"},
				String{Name: "second"},
			},
			input:     "Im first arg | im second arg",
			delimiter: " | ",
			expectErr: false,
			expectedArgsResults: map[string]any{
				"first":  "Im first arg",
				"second": "im second arg",
			},
		},
	}

	for _, c := range cases {
		t.Run(
			c.name, func(t *testing.T) {
				parser, err := NewParser(
					Opts{
						Args:          c.args,
						Input:         c.input,
						ArgsDelimiter: c.delimiter,
					},
				)
				if c.expectErr && err == nil {
					t.Errorf("expected error, got nil")
				}

				if err != nil {
					if !c.expectErr {
						t.Errorf("unexpected error: %v", err)
					} else {
						return
					}
				}

				if c.expectedArgsResults != nil && len(c.expectedArgsResults) > 0 {
					for k, v := range c.expectedArgsResults {
						arg := parser.Get(k)
						if v == nil && arg == nil {
							continue
						}

						if parser.Get(k).Value() != v {
							t.Errorf("expected %s to be %s, got %s", k, v, parser.Get(k).Value())
						}
					}
				}
			},
		)
	}
}
