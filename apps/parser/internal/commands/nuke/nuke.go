package nuke

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/xhit/go-str2duration/v2"
	"go.uber.org/zap"
)

const (
	nukePhraseArgName = "phrase"
	nukeTimeArgName   = "time"
	nukeRedisPrefix   = "channels:%s:nuked_messages"
)

var Command = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name: "nuke",
		Description: null.StringFrom(
			"Mass remove messages in chat by message content. Usage: !nuke 10m phrase, !nuke 10 phrase, !nuke 1h5m phrase",
		),
		RolesIDS: pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:   "MODERATION",
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name:     nukeTimeArgName,
			Optional: false,
		},
		command_arguments.VariadicString{
			Name: nukePhraseArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		duration, err := parseDuration(parseCtx.ArgsParser.Get(nukeTimeArgName).String())
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "invalid duration. Examples: !nuke 10m phrase, !nuke 10 phrase, !nuke 1h5m phrase",
			}
		}

		phrase := parseCtx.ArgsParser.Get(nukePhraseArgName).String()

		var messages []model.ChatMessage
		if err := parseCtx.Services.Gorm.
			Debug().
			Where("channel_id = ?", parseCtx.Channel.ID).
			Where("created_at > ?", time.Now().Add(-10*time.Minute)).
			Where("text ILIKE ?", "%"+phrase+"%").
			Joins("User").
			Joins(
				"User.Stats",
				parseCtx.Services.Gorm.Where(&model.UsersStats{ChannelID: parseCtx.Channel.ID}),
			).
			Where(
				`"User__Stats"."is_mod" = ? AND "User__Stats"."is_vip" = ? AND "User__Stats"."is_subscriber" = ?`,
				false, false, false,
			).
			Where(`"User"."id" != ?`, parseCtx.Channel.ID).
			Find(&messages).Error; err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get messages",
				Err:     err,
			}
		}

		if len(messages) == 0 {
			return nil, nil
		}

		handledMessagesIds, err := parseCtx.Services.Redis.SMembers(
			ctx,
			fmt.Sprintf(nukeRedisPrefix, parseCtx.Channel.ID),
		).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return nil, &types.CommandHandlerError{
				Message: "cannot get handled messages",
				Err:     err,
			}
		}

		mappedMessagesIDs := make([]string, 0, len(messages))
		for _, m := range messages {
			if !slices.Contains(handledMessagesIds, m.ID.String()) {
				mappedMessagesIDs = append(mappedMessagesIDs, m.ID.String())
			}
		}

		if duration <= 0 {
			if err := parseCtx.Services.Bus.Bots.DeleteMessage.Publish(
				bots.DeleteMessageRequest{
					ChannelId:   parseCtx.Channel.ID,
					MessageIds:  mappedMessagesIDs,
					ChannelName: &parseCtx.Channel.Name,
				},
			); err != nil {
				return nil, &types.CommandHandlerError{
					Message: "cannot delete messages",
					Err:     err,
				}
			}
		} else {
			var wg sync.WaitGroup

			for _, m := range messages {
				if slices.Contains(handledMessagesIds, m.ID.String()) {
					continue
				}

				wg.Add(1)
				m := m

				go func() {
					defer wg.Done()
					if err := parseCtx.Services.Bus.Bots.BanUser.Publish(
						bots.BanRequest{
							ChannelID:      parseCtx.Channel.ID,
							UserID:         m.UserID,
							Reason:         "nuked by twir",
							BanTime:        duration,
							IsModerator:    false,
							AddModAfterBan: false,
						},
					); err != nil {
						parseCtx.Services.Logger.Error("cannot ban user", zap.Error(err))
					}
				}()
			}

			wg.Wait()
		}

		parseCtx.Services.Redis.SAdd(
			ctx,
			fmt.Sprintf(nukeRedisPrefix, parseCtx.Channel.ID),
			mappedMessagesIDs,
		)
		parseCtx.Services.Redis.Expire(
			ctx,
			fmt.Sprintf(nukeRedisPrefix, parseCtx.Channel.ID),
			20*time.Minute,
		)

		return nil, nil
	},
}

func parseDuration(input string) (int, error) {
	asNumber, err := strconv.Atoi(input)
	if err == nil {
		return asNumber, nil
	}

	durationFromString, err := str2duration.ParseDuration(input)
	if durationFromString.Hours() >= 336 { // 2 weeks
		return 0, fmt.Errorf("duration of timeout cannot be longer than 2 weeks")
	}
	if err == nil {
		return int(durationFromString.Seconds()), nil
	}

	return 0, fmt.Errorf("cannot parse duration")
}
