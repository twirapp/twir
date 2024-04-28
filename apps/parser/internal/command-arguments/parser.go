package command_arguments

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Parser struct {
	store map[string]Arg
}

func (c *Parser) Get(name string) Arg {
	return c.store[name]
}

var ErrInvalidCommand = errors.New("[bug] command is invalid")
var ErrInvalidArg = errors.New("arg is invalid")

func NewParser(args []Arg, input string) (*Parser, error) {
	store := &Parser{}

	store.store = make(map[string]Arg)
	inputFields := strings.Fields(input)

	if len(args) == 0 {
		return store, nil
	}

	argsLen := len(args)
	if args[argsLen-1].IsOptional() {
		argsLen--
	}

	if len(inputFields) < argsLen {
		return store, ErrInvalidCommand
	}

	for i, arg := range args {
		if i >= len(inputFields) && arg.IsOptional() {
			break
		}

		switch typedArgument := arg.(type) {
		case VariadicString:
			value := strings.Join(inputFields[i:], " ")

			store.store[typedArgument.Name] = &VariadicString{
				Name:  typedArgument.Name,
				value: value,
			}
		case String:
			if len(typedArgument.OneOf) > 0 {
				if !slices.Contains(typedArgument.OneOf, inputFields[i]) {
					return store, ErrInvalidArg
				}
			}

			value := inputFields[i]

			store.store[typedArgument.Name] = &String{
				Name:  typedArgument.Name,
				value: value,
			}
			continue
		case Int:
			parsedInt, err := strconv.Atoi(inputFields[i])
			if err != nil {
				return store, ErrInvalidArg
			}

			if typedArgument.Max != nil && parsedInt > *typedArgument.Max {
				return store, ErrInvalidArg
			}

			if typedArgument.Min != nil && parsedInt < *typedArgument.Min {
				return store, ErrInvalidArg
			}

			store.store[typedArgument.Name] = &Int{
				Name:  typedArgument.Name,
				value: parsedInt,
			}
		}
	}

	return store, nil
}

func (c *Parser) BuildUsageString(args []Arg, cmdName string) string {
	usage := "!" + cmdName

	for _, arg := range args {
		switch typedArgument := arg.(type) {
		case VariadicString:
			usage += fmt.Sprintf(" [%s]", typedArgument.Name)
		case String:
			if len(typedArgument.OneOf) > 0 {
				usage += fmt.Sprintf(
					" <%s (%s)>",
					typedArgument.Name,
					strings.Join(typedArgument.OneOf, "|"),
				)
			} else {
				usage += fmt.Sprintf(" <%s>", typedArgument.Name)
			}
		case Int:
			if typedArgument.Min != nil && typedArgument.Max != nil {
				usage += fmt.Sprintf(
					" <%s (min %d, max %d)>",
					typedArgument.Name,
					*typedArgument.Min,
					*typedArgument.Max,
				)
			} else if typedArgument.Min != nil {
				usage += fmt.Sprintf(" <%s (min %d)>", typedArgument.Name, *typedArgument.Min)
			} else if typedArgument.Max != nil {
				usage += fmt.Sprintf(" <%s (max %d)>", typedArgument.Name, *typedArgument.Max)
			} else {
				usage += fmt.Sprintf(" <%s>", typedArgument.Name)
			}
		}
	}

	return usage
}
