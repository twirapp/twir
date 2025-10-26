package command_arguments

import (
	"context"
)

type Arg interface {
	isCommandArg()

	GetName() string
	String() string
	Int() int
	// IsOptional only latest argument can be optional
	IsOptional() bool
	GetHint(ctx context.Context) string
	Value() any
}
