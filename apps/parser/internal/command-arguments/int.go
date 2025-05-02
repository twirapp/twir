package command_arguments

type Int struct {
	value    int
	Name     string
	Min      *int
	Max      *int
	Optional bool
	Hint     string
}

var _ Arg = Int{}

func (Int) isCommandArg() {}

func (c Int) Int() int {
	return c.value
}

func (c Int) String() string {
	return ""
}

func (c Int) GetName() string {
	return c.Name
}

func (c Int) GetHint() string {
	if c.Hint == "" {
		return c.Name
	}

	return c.Hint
}

func (c Int) IsOptional() bool {
	return c.Optional
}

func (c Int) Value() any {
	return c.value
}
