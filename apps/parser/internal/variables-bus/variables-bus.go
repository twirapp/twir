package variables_bus

import (
	"cmp"
	"context"
	"slices"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/variables"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/parser"
)

type VariablesBus struct {
	bus  *buscore.Bus
	vars []parser.BuiltInVariable
}

func New(
	bus *buscore.Bus,
	variablesService *variables.Variables,
) *VariablesBus {
	b := &VariablesBus{
		bus:  bus,
		vars: make([]parser.BuiltInVariable, 0, len(variablesService.Store)),
	}

	for _, variable := range variablesService.Store {
		links := make([]parser.BuiltInVariableLink, 0, len(variable.Links))
		for _, l := range variable.Links {
			links = append(
				links, parser.BuiltInVariableLink{
					Name: l.Name,
					Href: l.Href,
				},
			)
		}

		b.vars = append(
			b.vars,
			parser.BuiltInVariable{
				Name:                variable.Name,
				Example:             lo.FromPtr(variable.Example),
				Description:         lo.FromPtr(variable.Description),
				Visible:             lo.FromPtr(variable.Visible),
				CanBeUsedInRegistry: variable.CanBeUsedInRegistry,
				Links:               links,
			},
		)
	}

	slices.SortFunc(
		b.vars,
		func(a, b parser.BuiltInVariable) int {
			return cmp.Compare(a.Name, b.Name)
		},
	)

	return b
}

func (c *VariablesBus) Subscribe() error {
	if err := c.bus.Parser.GetBuiltInVariables.SubscribeGroup(
		"parser",
		func(ctx context.Context, _ struct{}) ([]parser.BuiltInVariable, error) {
			return c.vars, nil
		},
	); err != nil {
		return err
	}

	return nil
}

func (c *VariablesBus) Unsubscribe() error {
	c.bus.Parser.GetBuiltInVariables.Unsubscribe()

	return nil
}
