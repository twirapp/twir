package types

import (
	"context"
)

type VariableHandlerResult struct {
	Result               string
	RepeatVariableResult *int
}

type VariableData struct {
	Params *string
	Key    string
}

type VariableHandler func(
	ctx context.Context,
	parseCtx *VariableParseContext,
	variableData *VariableData,
) (*VariableHandlerResult, error)

type Variable struct {
	Handler                  VariableHandler
	Description              *string
	Example                  *string
	Visible                  *bool
	Name                     string
	CommandsOnly             bool
	DisableInCustomVariables bool
	CanBeUsedInRegistry      bool
	NotCachable              bool
	Priority                 int // Higher number = higher priority, default 0
	Links                    []VariableLink
}

type VariableLink struct {
	Name string
	Href string
}

type VariableParseContext struct {
	*ParseContext
}
