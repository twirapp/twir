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
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/bus-core/bots"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"github.com/twirapp/twir/libs/repositories/chat_messages"
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
			Hint:     i18n.Get(locales.Translations.Commands.Nuke.Hints.NukeTimeArgName),
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
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Nuke.Errors.InvalidDuration),
			}
		}

		phrase := parseCtx.ArgsParser.Get(nukePhraseArgName).String()

		timeGte := time.Now().Add(-10 * time.Minute)

		messages, err := parseCtx.Services.ChatMessagesRepo.GetMany(
			ctx,
			chat_messages.GetManyInput{
				ChannelID: &parseCtx.Channel.ID,
				TextLike:  &phrase,
				Page:      0,
				PerPage:   1000,
				TimeGte:   &timeGte,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Errors.Generic.CannotGetMessage),
				Err:     err,
			}
		}
		if len(messages) == 0 {
			return nil, nil
		}

		usersIdsForCheck := make([]string, 0, len(messages))
		for _, m := range messages {
			usersIdsForCheck = append(usersIdsForCheck, m.UserID)
		}

		var usersStats []model.UsersStats
		if err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"userId" IN ? AND "channelId" = ?`, usersIdsForCheck, parseCtx.Channel.ID).
			Where(`"is_mod" = ? AND "is_vip" = ? AND "is_subscriber" = ?`, false, false, false).
			Where(`"userId" != ?`, parseCtx.Channel.ID).
			Find(&usersStats).Error; err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Nuke.Errors.CannotGetUsersStats),
				Err:     err,
			}
		}

		handledMessagesIds, err := parseCtx.Services.Redis.SMembers(
			ctx,
			fmt.Sprintf(nukeRedisPrefix, parseCtx.Channel.ID),
		).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Nuke.Errors.CannotGetHandeledMessages),
				Err:     err,
			}
		}

		mappedMessagesIDs := make([]string, 0, len(messages))
		for _, m := range messages {
			if slices.Contains(handledMessagesIds, m.ID.String()) {
				continue
			}

			if !slices.ContainsFunc(
				usersStats, func(stats model.UsersStats) bool {
					return m.UserID == stats.UserID
				},
			) {
				continue
			}

			mappedMessagesIDs = append(mappedMessagesIDs, m.ID.String())
		}

		if duration <= 0 {
			if err := parseCtx.Services.Bus.Bots.DeleteMessage.Publish(
				ctx,
				bots.DeleteMessageRequest{
					ChannelId:   parseCtx.Channel.ID,
					MessageIds:  mappedMessagesIDs,
					ChannelName: &parseCtx.Channel.Name,
				},
			); err != nil {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(ctx, locales.Translations.Commands.Nuke.Errors.CannotDeleteMessages),
					Err:     err,
				}
			}
		} else {
			var wg sync.WaitGroup

			var handledUsersIds []string

			for _, m := range messages {
				if slices.Contains(handledMessagesIds, m.ID.String()) {
					continue
				}

				if slices.Contains(handledUsersIds, m.UserID) {
					continue
				}
				handledUsersIds = append(handledUsersIds, m.UserID)

				wg.Add(1)
				m := m

				go func() {
					defer wg.Done()
					if err := parseCtx.Services.Bus.Bots.BanUser.Publish(
						ctx,
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
	if durationFromString.Hours() > 336 { // 2 weeks
		return 0, fmt.Errorf(i18n.Get(locales.Translations.Commands.Nuke.Errors.TimeoutDuration))
	}
	if err == nil {
		return int(durationFromString.Seconds()), nil
	}

	return 0, fmt.Errorf(i18n.Get(locales.Translations.Commands.Nuke.Errors.ParseDuration))
}
