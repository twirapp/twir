package spam

import (
	"context"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
)

const (
	spamCountArgName   = "count"
	spamMessageArgName = "message_or_command"
)

var Command = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:               "spam",
		Description:        null.StringFrom("Spam into chat. Example usage: !spam 5 Follow me on twitter or !spam 5 !tg"),
		RolesIDS:           pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:             "MODERATION",
		KeepResponsesOrder: false,
	},
	Args: []command_arguments.Arg{
		command_arguments.Int{
			Name: spamCountArgName,
			Min:  lo.ToPtr(1),
			Max:  lo.ToPtr(20),
		},
		command_arguments.VariadicString{
			Name: spamMessageArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{}

		count := parseCtx.ArgsParser.Get(spamCountArgName).Int()
		text := parseCtx.ArgsParser.Get(spamMessageArgName).String()

		var commandsPrefix string
		commandPrefixEntity, _ := parseCtx.Services.CommandsPrefixCache.Get(ctx, parseCtx.Channel.ID)
		if commandPrefixEntity != channelscommandsprefixmodel.Nil {
			commandsPrefix = commandPrefixEntity.Prefix
		} else {
			commandsPrefix = "!"
		}

		// if not command
		if !strings.HasPrefix(text, commandsPrefix) {
			for i := 0; i < count; i++ {
				result.Result = append(result.Result, text)
			}

			return result, nil
		}

		cmds, err := parseCtx.Services.CommandsCache.Get(ctx, parseCtx.Channel.ID)
		if err != nil {
			return nil, err
		}

		var foundCmd *model.ChannelsCommands

		for _, cmd := range cmds {
			if cmd.Name == strings.TrimPrefix(text, commandsPrefix) {
				foundCmd = &cmd
				break
			}

			for _, alias := range cmd.Aliases {
				if alias == strings.TrimPrefix(text, commandsPrefix) {
					foundCmd = &cmd
					break
				}
			}
		}

		if foundCmd == nil {
			return nil, nil
		}

		for i := 0; i < count; i++ {
			for _, r := range foundCmd.Responses {
				if !r.Text.Valid {
					continue
				}

				result.Result = append(result.Result, r.Text.String)
			}
		}

		return result, nil
	},
}
