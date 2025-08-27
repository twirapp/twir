package prefix

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	botssettings "github.com/twirapp/twir/libs/bus-core/bots-settings"
	model "github.com/twirapp/twir/libs/gomodels"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	"go.uber.org/zap"
)

const setPrefixArgName = "prefix"

var SetPrefix = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name: "prefix set",
		Description: null.StringFrom(
			"Set prefix for commands",
		),
		RolesIDS: pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:   "MODERATION",
		IsReply:  true,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: setPrefixArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		prefixArg := parseCtx.ArgsParser.Get(setPrefixArgName)
		if prefixArg == nil {
			return nil, &types.CommandHandlerError{
				Message: "prefix is required",
			}
		}

		if utf8.RuneCountInString(prefixArg.String()) > 10 {
			return nil, &types.CommandHandlerError{
				Message: "prefix cannot be longer than 10 characters",
			}
		}

		currentPrefix, err := parseCtx.Services.CommandsPrefixCache.Get(
			ctx,
			parseCtx.Channel.ID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get current prefix",
				Err:     err,
			}
		}

		var newPrefix channelscommandsprefixmodel.ChannelsCommandsPrefix

		newPrefix, err = parseCtx.Services.CommandsPrefixRepository.Update(
			ctx,
			currentPrefix.ID,
			channelscommandsprefixrepository.UpdateInput{
				Prefix: prefixArg.String(),
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot update prefix",
				Err:     err,
			}
		}

		go func() {
			if err = parseCtx.Services.Bus.BotsSettings.UpdatePrefix.Publish(
				ctx, botssettings.UpdatePrefixRequest{
					ID:        newPrefix.ID,
					ChannelID: newPrefix.ChannelID,
					Prefix:    newPrefix.Prefix,
					CreatedAt: newPrefix.CreatedAt,
					UpdatedAt: newPrefix.UpdatedAt,
				},
			); err != nil {
				parseCtx.Services.Logger.Error(
					"failed to publish channel command prefix update",
					zap.String("channel_id", newPrefix.ChannelID),
					zap.Any("error", err),
				)
			}
		}()

		return &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf("Prefix successfully updated to \"%s\"", prefixArg.String()),
			},
		}, nil
	},
}
