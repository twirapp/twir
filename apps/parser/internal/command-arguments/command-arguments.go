package command_arguments

type Arg interface {
	isCommandArg()

	GetName() string
	String() string
	Int() int
	// IsOptional only latest argument can be optional
	IsOptional() bool
	GetHint() string
	Value() any
}
