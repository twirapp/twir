package command_arguments

import (
	"context"
)

type String struct {
	value    string
	Name     string
	Optional bool
	OneOf    []string
	Hint     string
	HintFunc func(ctx context.Context) string
}

var _ Arg = String{}

func (String) isCommandArg() {}

func (c String) Int() int {
	return 0
}

func (c String) String() string {
	return c.value
}
func (c String) GetName() string {
	return c.Name
}

func (c String) GetHint(ctx context.Context) string {
	if c.HintFunc != nil {
		return c.HintFunc(ctx)
	}

	if c.Hint == "" {
		return c.Name
	}

	return c.Hint
}

func (c String) IsOptional() bool {
	return c.Optional
}

func (c String) Value() any {
	return c.value
}
