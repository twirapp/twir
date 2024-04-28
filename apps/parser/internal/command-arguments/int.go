package command_arguments

type Int struct {
	value    int
	Name     string
	Min      *int
	Max      *int
	Optional bool
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

func (c Int) IsOptional() bool {
	return c.Optional
}
