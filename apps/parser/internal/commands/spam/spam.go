package spam

import (
	"context"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
)

const (
	spamCountArgName   = "count"
	spamMessageArgName = "message"
)

var Command = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "spam",
		Description: null.StringFrom("Spam into chat. Example usage: <b>!spam 5 Follow me on twitter"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MODERATION",
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

		for i := 0; i < count; i++ {
			result.Result = append(result.Result, text)
		}

		return result, nil
	},
}
