package command_arguments

import (
	"context"
)

type VariadicString struct {
	value    string
	Name     string
	Optional bool
	Hint     string
	HintFunc func(ctx context.Context) string
}

var _ Arg = VariadicString{}

func (VariadicString) isCommandArg() {}

func (c VariadicString) Int() int {
	return 0
}

func (c VariadicString) String() string {
	return c.value
}

func (c VariadicString) GetName() string {
	return c.Name
}

func (c VariadicString) GetHint(ctx context.Context) string {
	if c.HintFunc != nil {
		return c.HintFunc(ctx)
	}

	if c.Hint == "" {
		return c.Name
	}

	return c.Hint
}

func (c VariadicString) IsOptional() bool {
	return c.Optional
}

func (c VariadicString) Value() any {
	return c.value
}
