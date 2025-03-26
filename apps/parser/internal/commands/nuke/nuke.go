package nuke

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"go.uber.org/zap"
)

const (
	nukePhraseArgName = "phrase"
	nukeTimeArgName   = "time"
	nukeRedisPrefix   = "channels:%s:nuked_messages:%s"
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
			Where("channel_id = ?", parseCtx.Channel.ID).
			Where("created_at > ?", time.Now().Add(-10*time.Minute)).
			Where("text ILIKE ?", "%"+phrase+"%").
			Find(&messages).Error; err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get messages",
				Err:     err,
			}
		}

		if len(messages) == 0 {
			return nil, nil
		}

		var handledMessagesIds []string
		iter := parseCtx.Services.Redis.
			Scan(
				ctx,
				0,
				fmt.Sprintf(nukeRedisPrefix, parseCtx.Channel.ID, "*"),
				0,
			).
			Iterator()
		for iter.Next(ctx) {
			handledMessagesIds = append(
				handledMessagesIds,
				strings.Replace(iter.Val(), fmt.Sprintf(nukeRedisPrefix, parseCtx.Channel.ID, ""), "", 1),
			)
		}
		if err := iter.Err(); err != nil {
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

		parseCtx.Services.Redis.Pipelined(
			ctx,
			func(pipe redis.Pipeliner) error {
				for _, msgID := range mappedMessagesIDs {
					if err := pipe.Set(
						ctx,
						fmt.Sprintf(nukeRedisPrefix, parseCtx.Channel.ID, msgID),
						true,
						20*time.Minute,
					).Err(); err != nil {
						return err
					}
				}

				return nil
			},
		)

		return nil, nil
	},
}

func parseDuration(input string) (int, error) {
	asNumber, err := strconv.Atoi(input)
	if err == nil {
		return asNumber, nil
	}

	duration, err := time.ParseDuration(input)
	if err == nil {
		return int(duration.Seconds()), nil
	}

	return 0, fmt.Errorf("cannot parse duration")
}
