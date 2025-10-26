package user

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
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
			Where(`"channelId" = ? AND "orderedById" = ?`, parseCtx.Channel.ID, targetUserId).
			Count(&count).
			Error

		if err != nil {
			result.Result = "0"
			return result, nil
		}

		result.Result = i18n.GetCtx(
			ctx,
			locales.Translations.Variables.User.Info.Songs.
				SetVars(locales.KeysVariablesUserInfoSongsVars{UserSongs: count}),
		)

		return result, nil
	},
}
