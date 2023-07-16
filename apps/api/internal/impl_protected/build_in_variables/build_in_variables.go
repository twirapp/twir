package build_in_variables

import (
	"context"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	"github.com/satont/twir/libs/grpc/generated/api/built_in_variables"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BuildInVariables struct {
	*impl_deps.Deps
}

func (c *BuildInVariables) BuiltInVariablesGetAll(
	ctx context.Context,
	req *emptypb.Empty,
) (*built_in_variables.GetAllResponse, error) {
	variables, err := c.Deps.Grpc.Parser.GetDefaultVariables(ctx, req)
	if err != nil {
		return nil, err
	}

	return &built_in_variables.GetAllResponse{
		Variables: lo.Map(
			variables.List,
			func(item *parser.GetVariablesResponse_Variable, index int) *built_in_variables.Variable {
				return &built_in_variables.Variable{
					Name:        item.Name,
					Example:     item.Example,
					Description: item.Description,
					Visible:     item.Visible,
				}
			}),
	}, nil
}
