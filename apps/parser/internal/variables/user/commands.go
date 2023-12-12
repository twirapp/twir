package user

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
)

var Commands = &types.Variable{
	Name:         "user.commands",
	Description:  lo.ToPtr("User used commands count"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		var count int64
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND "userId" = ?`, parseCtx.Channel.ID, parseCtx.Sender.ID).
			Model(&model.ChannelsCommandsUsages{}).
			Count(&count).
			Error

		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)

			result.Result = "internal error"
			return result, nil
		}

		result.Result = fmt.Sprint(count)

		return result, nil
	},
}
