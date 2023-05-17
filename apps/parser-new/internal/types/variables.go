package types

import (
	"context"
)

type VariableHandlerResult struct {
	Result string
}

type VariableData struct {
	Key    string
	Params *string
}

type VariableHandler func(
	ctx context.Context,
	parseCtx *VariableParseContext,
	variableData *VariableData,
) (*VariableHandlerResult, error)

type Variable struct {
	Name         string
	Handler      VariableHandler
	Description  *string
	Example      *string
	CommandsOnly *bool
	Visible      *bool
}

type VariableParseContext struct {
	*ParseContext
}
