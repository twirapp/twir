package games

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/bus-core/bots"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"gorm.io/gorm"
)

var Seppuku = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "seppuku",
		Description: null.StringFrom("Seppuku, is a form of Japanese ritualistic suicide by disembowelment."),
		Module:      "GAMES",
		IsReply:     true,
		Visible:     true,
		RolesIDS:    pq.StringArray{},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		entity := model.ChannelGamesSeppuku{}
		if err := parseCtx.Services.Gorm.WithContext(ctx).Where(
			`"channel_id" = ?`,
			parseCtx.Channel.ID,
		).First(&entity).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil
			}

			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Games.Errors.SeppukuCannotFindSettings),
				Err:     err,
			}
		}

		if !entity.Enabled {
			return &types.CommandsHandlerResult{
				Result: []string{},
			}, nil
		}

		if !entity.TimeoutModerators && slices.Contains(parseCtx.Sender.Badges, "MODERATOR") {
			return &types.CommandsHandlerResult{
				Result: []string{},
			}, nil
		}

		var message string
		if slices.Contains(parseCtx.Sender.Badges, "MODERATOR") {
			message = entity.MessageModerators
		} else {
			message = entity.Message
		}
		message = strings.ReplaceAll(message, "{sender}", parseCtx.Sender.DisplayName)

		if parseCtx.Sender.ID == parseCtx.Channel.ID {
			return &types.CommandsHandlerResult{
				Result: []string{message},
			}, nil
		}

		isModerator := slices.Contains(parseCtx.Sender.Badges, "MODERATOR")
		if !entity.TimeoutModerators && isModerator {
			return &types.CommandsHandlerResult{
				Result: []string{message},
			}, nil
		}

		if err := parseCtx.Services.Bus.Bots.BanUser.Publish(
			ctx,
			bots.BanRequest{
				ChannelID:      parseCtx.Channel.ID,
				UserID:         parseCtx.Sender.ID,
				BanTime:        entity.TimeoutSeconds,
				Reason:         message,
				IsModerator:    isModerator,
				AddModAfterBan: true,
			},
		); err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Errors.Generic.CannotBanUser),
				Err:     err,
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{message},
		}, nil
	},
}
