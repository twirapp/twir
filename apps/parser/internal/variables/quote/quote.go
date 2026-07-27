package quote

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	"gorm.io/gorm"
)

var Quote = &types.Variable{
	Name:        "quote",
	Description: lo.ToPtr("Random quote or quote by number"),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}
		quote := model.ChannelsQuotes{}
		db := parseCtx.Services.Gorm.WithContext(ctx).Where(`"channelId" = ?`, parseCtx.Channel.DBChannelID)
		params := strings.TrimPrefix(strings.TrimSpace(*variableData.Params), "#")

		var err error
		if params == "" {
			err = db.Order("RANDOM()").First(&quote).Error
		} else {
			number, parseErr := strconv.Atoi(params)
			if parseErr != nil || number < 1 {
				return result, nil
			}
			err = db.Where("number = ?", number).First(&quote).Error
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, nil
		}
		if err != nil {
			return nil, err
		}

		result.Result = quote.Text
		return result, nil
	},
}
