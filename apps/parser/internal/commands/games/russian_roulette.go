package games

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/twitch"
	"golang.org/x/exp/rand"
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
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: []string{},
		}

		entity := model.ChannelModulesSettings{}
		if err := parseCtx.Services.Gorm.WithContext(ctx).Where(
			`"channelId" = ? and "userId" is null and "type" = 'russian_roulette'`,
			parseCtx.Channel.ID,
		).First(&entity).Error; err != nil {
			return result
		}

		var parsedSettings model.RussianRouletteSetting
		if err := json.Unmarshal(entity.Settings, &parsedSettings); err != nil {
			return result
		}

		if !parsedSettings.Enabled {
			return result
		}

		initMessage := strings.ReplaceAll(
			parsedSettings.InitMessage,
			"{sender}",
			parseCtx.Sender.DisplayName,
		)
		surviveMessage := strings.ReplaceAll(
			parsedSettings.SurviveMessage,
			"{sender}",
			parseCtx.Sender.DisplayName,
		)
		deathMessage := strings.ReplaceAll(
			parsedSettings.DeathMessage,
			"{sender}",
			parseCtx.Sender.DisplayName,
		)

		_, err := parseCtx.Services.GrpcClients.Bots.SendMessage(
			ctx, &bots.SendMessageRequest{
				ChannelId:      parseCtx.Channel.ID,
				ChannelName:    &parseCtx.Channel.Name,
				Message:        initMessage,
				SkipRateLimits: true,
			},
		)
		if err != nil {
			return result
		}

		if parsedSettings.DecisionSeconds > 0 {
			time.Sleep(time.Duration(parsedSettings.DecisionSeconds) * time.Second)
		}

		if slices.Contains(parseCtx.Sender.Badges, "BROADCASTER") {
			result.Result = []string{surviveMessage}
			return result
		}

		twitchClient, err := twitch.NewUserClient(
			parseCtx.Channel.ID,
			*parseCtx.Services.Config,
			parseCtx.Services.GrpcClients.Tokens,
		)
		if err != nil {
			return result
		}

		randomized := rand.Intn(8)
		if randomized > parsedSettings.ChargedBullets {
			result.Result = []string{surviveMessage}
			return result
		} else {
			isModerator := slices.Contains(parseCtx.Sender.Badges, "MODERATOR")
			if parsedSettings.CanBeUsedByModerators && isModerator && parsedSettings.TimeoutSeconds > 0 {
				_, err = twitchClient.RemoveChannelModerator(
					&helix.RemoveChannelModeratorParams{
						UserID:        parseCtx.Sender.ID,
						BroadcasterID: parseCtx.Channel.ID,
					},
				)
				if err != nil {
					result.Result = []string{"internal error when trying to remove mod"}
					return result
				}

				go func() {
					time.Sleep(time.Duration(parsedSettings.TimeoutSeconds+2) * time.Second)

					_, err = twitchClient.AddChannelModerator(
						&helix.AddChannelModeratorParams{
							UserID:        parseCtx.Sender.ID,
							BroadcasterID: parseCtx.Channel.ID,
						},
					)
					if err != nil {
						return
					}
				}()
			}

			if parsedSettings.TimeoutSeconds > 0 &&
				(!isModerator || (isModerator && parsedSettings.CanBeUsedByModerators)) {
				_, err = twitchClient.BanUser(
					&helix.BanUserParams{
						BroadcasterID: parseCtx.Channel.ID,
						ModeratorId:   parseCtx.Channel.ID,
						Body: helix.BanUserRequestBody{
							Duration: parsedSettings.TimeoutSeconds,
							Reason:   deathMessage,
							UserId:   parseCtx.Sender.ID,
						},
					},
				)
				if err != nil {
					result.Result = []string{"internal error when trying to ban user"}
					return result
				}
			}

			_, err = parseCtx.Services.GrpcClients.Bots.SendMessage(
				ctx,
				&bots.SendMessageRequest{
					ChannelId:      parseCtx.Channel.ID,
					ChannelName:    &parseCtx.Channel.Name,
					Message:        deathMessage,
					SkipRateLimits: true,
				},
			)
			if err != nil {
				return result
			}

			result.Result = []string{}
			return result
		}
	},
}
