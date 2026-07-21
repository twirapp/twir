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
	Description:  new("How many songs user requested"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		targetUserId := lo.
			IfF(
				len(parseCtx.Mentions) > 0, func() string {
					return parseCtx.Mentions[0].UserID
				},
			).
			Else(parseCtx.Sender.ID)

		user, err := parseCtx.Services.UsersRepo.GetByPlatformID(ctx, parseCtx.Platform, targetUserId)
		if err != nil {
			return nil, err
		}

		var count int64
		err = parseCtx.Services.Gorm.
			WithContext(ctx).
			Model(&model.RequestedSong{}).
			Where(
				`"channelId" = ?::uuid AND "orderedById" = ?`,
				parseCtx.Channel.DBChannelID,
				user.ID.String(),
			).
			Count(&count).
			Error

		if err != nil {
			result.Result = "0"
			return result, nil
		}

		formattedCount := strconv.Itoa(int(count))
		result.Result = formattedCount

		return result, nil
	},
}
