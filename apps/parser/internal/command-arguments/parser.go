package command_arguments

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Parser struct {
	store     map[string]Arg
	delimiter string
}

func (c *Parser) Get(name string) Arg {
	return c.store[name]
}

func (c *Parser) IsExists(name string) bool {
	_, exists := c.store[name]
	return exists
}

var ErrInvalidCommand = errors.New("[bug] command is invalid")
var ErrInvalidArg = errors.New("arg is invalid")

type Opts struct {
	Args          []Arg
	Input         string
	ArgsDelimiter string
}

func NewParser(opts Opts) (*Parser, error) {
	p := &Parser{
		store:     make(map[string]Arg),
		delimiter: " ",
	}
	if opts.ArgsDelimiter != "" {
		p.delimiter = opts.ArgsDelimiter
	}

	inputFields := strings.Split(opts.Input, p.delimiter)
	inputFields = slices.DeleteFunc(
		inputFields, func(s string) bool {
			return s == ""
		},
	)

	if len(opts.Args) == 0 {
		return p, nil
	}

	argsLen := len(opts.Args)
	if opts.Args[argsLen-1].IsOptional() {
		argsLen--
	}

	if len(inputFields) < argsLen {
		return p, ErrInvalidCommand
	}

	for i, arg := range opts.Args {
		if i >= len(inputFields) && arg.IsOptional() {
			break
		}

		switch typedArgument := arg.(type) {
		case VariadicString:
			value := strings.Join(inputFields[i:], " ")

			p.store[typedArgument.Name] = &VariadicString{
				Name:  typedArgument.Name,
				value: value,
			}
		case String:
			if len(typedArgument.OneOf) > 0 {
				if !slices.Contains(typedArgument.OneOf, inputFields[i]) {
					return p, ErrInvalidArg
				}
			}

			value := inputFields[i]

			p.store[typedArgument.Name] = &String{
				Name:  typedArgument.Name,
				value: value,
			}
			continue
		case Int:
			parsedInt, err := strconv.Atoi(inputFields[i])
			if err != nil {
				return p, ErrInvalidArg
			}

			if typedArgument.Max != nil && parsedInt > *typedArgument.Max {
				return p, ErrInvalidArg
			}

			if typedArgument.Min != nil && parsedInt < *typedArgument.Min {
				return p, ErrInvalidArg
			}

			p.store[typedArgument.Name] = &Int{
				Name:  typedArgument.Name,
				value: parsedInt,
			}
		}
	}

	return p, nil
}

func (c *Parser) BuildUsageString(args []Arg, cmdName string) string {
	var usageBuidler strings.Builder

	usageBuidler.WriteString("!")
	usageBuidler.WriteString(cmdName)
	usageBuidler.WriteString(" ")

	for idx, arg := range args {
		if idx != 0 {
			usageBuidler.WriteString(c.delimiter)
		}

		switch typedArgument := arg.(type) {
		case VariadicString:
			usageBuidler.WriteString(c.buildUsagePrefixAndSuffix(typedArgument.GetHint()))
		case String:
			if len(typedArgument.OneOf) > 0 {
				usageBuidler.WriteString(
					c.buildUsagePrefixAndSuffix(
						fmt.Sprintf("%s %s", typedArgument.Name, typedArgument.GetHint()),
					),
				)
			} else {
				usageBuidler.WriteString(c.buildUsagePrefixAndSuffix(typedArgument.GetHint()))
			}
		case Int:
			if typedArgument.Min != nil && typedArgument.Max != nil {
				usageBuidler.WriteString(
					c.buildUsagePrefixAndSuffix(
						fmt.Sprintf(
							"%s (min %d, max %d)",
							typedArgument.GetHint(),
							*typedArgument.Min,
							*typedArgument.Max,
						),
					),
				)
			} else if typedArgument.Min != nil {
				usageBuidler.WriteString(
					c.buildUsagePrefixAndSuffix(
						fmt.Sprintf(
							"%s (min %d)",
							typedArgument.GetHint(),
							*typedArgument.Min,
						),
					),
				)
			} else if typedArgument.Max != nil {
				usageBuidler.WriteString(
					c.buildUsagePrefixAndSuffix(
						fmt.Sprintf(
							"%s (max %d)",
							typedArgument.GetHint(),
							*typedArgument.Max,
						),
					),
				)
			} else {
				usageBuidler.WriteString(c.buildUsagePrefixAndSuffix(typedArgument.GetHint()))
			}
		}
	}

	return usageBuidler.String()
}

func (c *Parser) buildUsagePrefixAndSuffix(input string) string {
	prefix := "<"
	suffix := ">"
	if c.delimiter != " " {
		prefix = ""
		suffix = ""
	}

	return prefix + input + suffix
}
