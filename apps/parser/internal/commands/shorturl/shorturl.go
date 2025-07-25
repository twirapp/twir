package shorturl

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
)

const (
	urlArgName = "url"
)

var Command = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "shorturl",
		Description: null.StringFrom("Create short url"),
		RolesIDS:    pq.StringArray{},
		Module:      "UTILITY",
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: urlArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{}

		urlArg := parseCtx.ArgsParser.Get(urlArgName).String()

		link, err := parseCtx.Services.ShortUrlServices.FindOrCreate(
			ctx,
			urlArg,
			parseCtx.Sender.ID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("cannot create short url: %s", err),
				Err:     err,
			}
		}

		result.Result = []string{
			fmt.Sprintf("✅ %s", link.Short),
		}
		return result, nil
	},
}
