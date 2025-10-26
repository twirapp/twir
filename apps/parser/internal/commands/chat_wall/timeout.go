package chat_wall

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	chatwallservice "github.com/twirapp/twir/apps/parser/internal/services/chat_wall"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	chatwallmodel "github.com/twirapp/twir/libs/repositories/chat_wall/model"
	"github.com/xhit/go-str2duration/v2"
)

const timeoutDurationArgName = "duration"
const timeoutPhraseArgName = "phrase"

var Timeout = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "chat wall timeout",
		Description: null.StringFrom("Creates chat wall instance for 10 minutes, which will timeout all messages with specific phrase"),
		RolesIDS: pq.StringArray{
			model.ChannelRoleTypeModerator.String(),
		},
		Module:  "MODERATION",
		Visible: true,
		IsReply: true,
	},
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: timeoutDurationArgName,
			HintFunc: func(ctx context.Context) string {
				return i18n.GetCtx(ctx, locales.Translations.Commands.ChatWall.Hints.TimeoutDurationArgName)
			},
		},
		command_arguments.VariadicString{
			Name: timeoutPhraseArgName,
			HintFunc: func(ctx context.Context) string {
				return i18n.GetCtx(ctx, locales.Translations.Commands.ChatWall.Hints.TimeoutPhraseArgName)
			},
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		phrase := parseCtx.ArgsParser.Get(timeoutPhraseArgName).String()
		parsedDuration, err := parseDuration(parseCtx.ArgsParser.Get(timeoutDurationArgName).String())
		if err != nil || parsedDuration <= 0 {
			return nil, &types.CommandHandlerError{
				Message: i18n.Get(locales.Translations.Commands.ChatWall.Errors.InvalidDuration),
			}
		}

		timeoutDuration := time.Duration(parsedDuration) * time.Second

		wall, err := parseCtx.Services.ChatWallService.Create(
			ctx,
			chatwallservice.CreateInput{
				ChannelID:       parseCtx.Channel.ID,
				Phrase:          phrase,
				Enabled:         true,
				Action:          chatwallmodel.ChatWallActionTimeout,
				Duration:        10 * time.Minute,
				TimeoutDuration: &timeoutDuration,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: err.Error(),
				Err:     err,
			}
		}

		if err := parseCtx.Services.ChatWallService.HandlePastMessages(
			ctx,
			wall,
			chatwallservice.HandlePastMessagesInput{
				ChannelID:       parseCtx.Channel.ID,
				Phrase:          phrase,
				Action:          chatwallmodel.ChatWallActionTimeout,
				TimeoutDuration: &timeoutDuration,
			},
		); err != nil {
			return nil, &types.CommandHandlerError{
				Message: err.Error(),
				Err:     err,
			}
		}

		result := &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.ChatWall.Start.ChatWallStart.
						SetVars(locales.KeysCommandsChatWallStartChatWallStartVars{ChatWallPhrase: phrase}),
				),
			},
		}

		return result, nil
	},
}

func parseDuration(input string) (int, error) {
	asNumber, err := strconv.Atoi(input)
	if err == nil {
		return asNumber, nil
	}

	durationFromString, err := str2duration.ParseDuration(input)
	if durationFromString.Hours() > 336 { // 2 weeks
		return 0, fmt.Errorf(i18n.Get(locales.Translations.Commands.ChatWall.Errors.LongDurationTimeout))
	}
	if err == nil {
		return int(durationFromString.Seconds()), nil
	}

	return 0, fmt.Errorf(i18n.Get(locales.Translations.Commands.ChatWall.Errors.DurationCannotParse))
}
