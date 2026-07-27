package quote

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	quotesmodel "github.com/twirapp/twir/libs/repositories/quotes/model"
)

var Quote = &types.Variable{
	Name:        "quote",
	Description: lo.ToPtr("Random quote or quote by number"),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		channelUUID, err := uuid.Parse(parseCtx.Channel.DBChannelID)
		if err != nil {
			return nil, err
		}

		params := strings.TrimPrefix(strings.TrimSpace(*variableData.Params), "#")

		var quote quotesmodel.Quote
		if params == "" {
			quote, err = parseCtx.Services.QuotesRepo.GetRandomByChannelID(ctx, channelUUID)
		} else {
			number, parseErr := strconv.Atoi(params)
			if parseErr != nil || number < 1 {
				return result, nil
			}
			quote, err = parseCtx.Services.QuotesRepo.GetByChannelIDAndNumber(ctx, channelUUID, number)
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return result, nil
		}
		if err != nil {
			return nil, err
		}

		result.Result = quote.Text
		return result, nil
	},
}
