package types

import (
	"context"
)

type VariableHandlerResult struct {
	Result string
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
}

type VariableParseContext struct {
	*ParseContext
}
