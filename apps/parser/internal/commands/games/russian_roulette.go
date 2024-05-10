package games

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/guregu/null"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/task-queue"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/libs/bus-core/bots"
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
				Message: "cannot get roulette settings from db",
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
				Message: "cannot send initial message",
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

		twitchClient, err := twitch.NewUserClient(
			parseCtx.Channel.ID,
			*parseCtx.Services.Config,
			parseCtx.Services.GrpcClients.Tokens,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create broadcaster twitch client",
				Err:     err,
			}
		}

		randomized := rand.Intn(entity.TumberSize + 1)
		if randomized > entity.ChargedBullets {
			result.Result = []string{surviveMessage}
			return result, nil
		} else {
			parseCtx.Services.Bus.Bots.SendMessage.Publish(
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
					Message: "cannot send death message",
					Err:     err,
				}
			}

			isModerator := slices.Contains(parseCtx.Sender.Badges, "MODERATOR")
			if entity.CanBeUsedByModerators && isModerator && entity.TimeoutSeconds > 0 {
				err = parseCtx.Services.TaskDistributor.DistributeModUser(
					ctx,
					&task_queue.TaskModUserPayload{
						ChannelID: parseCtx.Channel.ID,
						UserID:    parseCtx.Sender.ID,
					},
					asynq.ProcessIn(time.Duration(entity.TimeoutSeconds+2)*time.Second),
				)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: "cannot distribute mod user",
						Err:     err,
					}
				}

				_, err = twitchClient.RemoveChannelModerator(
					&helix.RemoveChannelModeratorParams{
						UserID:        parseCtx.Sender.ID,
						BroadcasterID: parseCtx.Channel.ID,
					},
				)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: "cannot remove moderator",
						Err:     err,
					}
				}
			}

			if entity.TimeoutSeconds > 0 &&
				(!isModerator || (isModerator && entity.CanBeUsedByModerators)) {
				_, err = twitchClient.BanUser(
					&helix.BanUserParams{
						BroadcasterID: parseCtx.Channel.ID,
						ModeratorId:   parseCtx.Channel.ID,
						Body: helix.BanUserRequestBody{
							Duration: entity.TimeoutSeconds,
							Reason:   deathMessage,
							UserId:   parseCtx.Sender.ID,
						},
					},
				)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: "cannot ban user",
						Err:     err,
					}
				}
			}

			result.Result = []string{}
			return result, nil
		}
	},
}
