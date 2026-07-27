package variables

import (
	"context"
	"testing"

	"github.com/twirapp/twir/apps/parser/internal/types"
)

func TestParseVariablesInText_whenVariableHasSpaceArgument(t *testing.T) {
	variables := &Variables{
		Store: map[string]*types.Variable{
			"quote": {
				Handler: func(
					ctx context.Context,
					parseCtx *types.VariableParseContext,
					variableData *types.VariableData,
				) (*types.VariableHandlerResult, error) {
					return &types.VariableHandlerResult{Result: *variableData.Params}, nil
				},
			},
		},
	}

	result := variables.ParseVariablesInText(context.Background(), &types.ParseContext{}, "$(quote 42)")
	if len(result) != 1 || result[0] != "42" {
		t.Fatalf("expected quote argument to be resolved, got %#v", result)
	}
}
