package chat_wall

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	chatwallservice "github.com/satont/twir/apps/parser/internal/services/chat_wall"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
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
			Hint: "time, examples: 10m, 10, 1h5m",
		},
		command_arguments.VariadicString{
			Name: timeoutPhraseArgName,
			Hint: "phrase to ban",
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
				Message: "invalid duration. Cannot be longer 2w Examples: 10m, 10, 1h5m",
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
				fmt.Sprintf(
					`✅ Chat wall started for 10 minutes, you can stop it with !chat wall stop "%s"`,
					phrase,
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
		return 0, fmt.Errorf("duration of timeout cannot be longer than 2 weeks")
	}
	if err == nil {
		return int(durationFromString.Seconds()), nil
	}

	return 0, fmt.Errorf("cannot parse duration")
}
