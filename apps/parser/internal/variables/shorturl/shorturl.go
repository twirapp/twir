package shorturl

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var Variable = &types.Variable{
	Name:         "shorturl",
	Description:  lo.ToPtr("Create short url from your link"),
	Example:      lo.ToPtr("shorturl|https://example.com"),
	CommandsOnly: false,
	NotCachable:  true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		if variableData.Params == nil {
			return nil, &types.CommandHandlerError{
				Message: "url is required",
			}
		}

		link, err := parseCtx.Services.ShortUrlServices.FindOrCreate(
			ctx,
			*variableData.Params,
			parseCtx.Sender.ID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("cannot create short url: %s", err),
				Err:     err,
			}
		}

		result := types.VariableHandlerResult{Result: link.Short}

		return &result, nil
	},
}
