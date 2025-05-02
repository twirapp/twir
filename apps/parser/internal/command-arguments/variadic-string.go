package command_arguments

type VariadicString struct {
	value    string
	Name     string
	Optional bool
	Hint     string
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

func (c VariadicString) GetHint() string {
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
