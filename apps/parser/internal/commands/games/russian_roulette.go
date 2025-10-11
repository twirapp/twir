package games

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/bus-core/bots"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

var RussianRoulette = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "roulette",
		Description: null.StringFrom("Test your luck!"),
		Module:      "GAMES",
		IsReply:     true,
		Visible:     true,
		RolesIDS:    pq.StringArray{},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{
			Result: []string{},
		}

		entity := model.ChannelGamesRussianRoulette{}
		if err := parseCtx.Services.Gorm.WithContext(ctx).Where(
			`"channel_id" = ?`,
			parseCtx.Channel.ID,
		).First(&entity).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return result, nil
			}

			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Games.Errors.RouletteCannotGetWithSettings),
				Err:     err,
			}
		}

		if !entity.Enabled {
			return result, nil
		}

		initMessage := strings.ReplaceAll(
			entity.InitMessage,
			"{sender}",
			parseCtx.Sender.DisplayName,
		)
		surviveMessage := strings.ReplaceAll(
			entity.SurviveMessage,
			"{sender}",
			parseCtx.Sender.DisplayName,
		)
		deathMessage := strings.ReplaceAll(
			entity.DeathMessage,
			"{sender}",
			parseCtx.Sender.DisplayName,
		)

		replyTo := lo.IfF(
			parseCtx.Command.IsReply, func() string {
				return parseCtx.MessageId
			},
		).Else("")

		err := parseCtx.Services.Bus.Bots.SendMessage.Publish(
			ctx,
			bots.SendMessageRequest{
				ChannelId:      parseCtx.Channel.ID,
				ChannelName:    &parseCtx.Channel.Name,
				Message:        initMessage,
				SkipRateLimits: true,
				ReplyTo:        replyTo,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Games.Errors.RouletteCannotSendInitialMessage),
				Err:     err,
			}
		}

		if entity.DecisionSeconds > 0 {
			time.Sleep(time.Duration(entity.DecisionSeconds) * time.Second)
		}

		if slices.Contains(parseCtx.Sender.Badges, "BROADCASTER") {
			result.Result = []string{surviveMessage}
			return result, nil
		}

		randomized := rand.Intn(entity.TumberSize + 1)
		if randomized > entity.ChargedBullets {
			result.Result = []string{surviveMessage}
			return result, nil
		} else {
			parseCtx.Services.Bus.Bots.SendMessage.Publish(
				ctx,
				bots.SendMessageRequest{
					ChannelId:      parseCtx.Channel.ID,
					ChannelName:    &parseCtx.Channel.Name,
					Message:        deathMessage,
					SkipRateLimits: true,
					ReplyTo:        replyTo,
				},
			)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(ctx, locales.Translations.Commands.Games.Errors.RouletteCannotSendDeathMessage),
					Err:     err,
				}
			}

			isModerator := slices.Contains(parseCtx.Sender.Badges, "MODERATOR")

			if entity.TimeoutSeconds > 0 {
				if isModerator && !entity.CanBeUsedByModerators {
					return result, nil
				}

				err = parseCtx.Services.Bus.Bots.BanUser.Publish(
					ctx,
					bots.BanRequest{
						ChannelID:      parseCtx.Channel.ID,
						UserID:         parseCtx.Sender.ID,
						Reason:         deathMessage,
						BanTime:        entity.TimeoutSeconds,
						AddModAfterBan: true,
						IsModerator:    isModerator,
					},
				)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: i18n.GetCtx(ctx, locales.Translations.Errors.Generic.CannotBanUser),
						Err:     err,
					}
				}
			}

			result.Result = []string{}
			return result, nil
		}
	},
}
