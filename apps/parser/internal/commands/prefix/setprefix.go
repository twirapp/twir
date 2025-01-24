package prefix

import (
	"context"
	"errors"
	"unicode/utf8"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
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

		currentPrefix, err := parseCtx.Services.CommandsPrefixRepository.GetByChannelID(
			ctx,
			parseCtx.Channel.ID,
		)
		if err != nil && !errors.Is(err, channelscommandsprefixrepository.ErrNotFound) {
			return nil, &types.CommandHandlerError{
				Message: "cannot get current prefix",
				Err:     err,
			}
		}

		if currentPrefix == channelscommandsprefixmodel.Nil {
			_, err = parseCtx.Services.CommandsPrefixRepository.Create(
				ctx,
				channelscommandsprefixrepository.CreateInput{
					ChannelID: parseCtx.Channel.ID,
					Prefix:    prefixArg.String(),
				},
			)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: "cannot create prefix",
					Err:     err,
				}
			}
		} else {
			_, err = parseCtx.Services.CommandsPrefixRepository.Update(
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
		}

		return &types.CommandsHandlerResult{Result: []string{"Prefix updated"}}, nil
	},
}
