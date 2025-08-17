package user

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
)

var SongsRequested = &types.Variable{
	Name:         "user.songs.requested.count",
	Description:  lo.ToPtr("How many songs user requested"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		targetUserId := lo.
			IfF(
				len(parseCtx.Mentions) > 0, func() string {
					return parseCtx.Mentions[0].UserId
				},
			).
			Else(parseCtx.Sender.ID)
		var count int64
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Model(&model.RequestedSong{}).
			Where(`"channelId" = ? AND "orderedById" = ?"`, parseCtx.Channel.ID, targetUserId).
			Count(&count).
			Error

		if err != nil {
			result.Result = "0"
			return result, nil
		}

		result.Result = strconv.FormatInt(count, 10)

		return result, nil
	},
}
