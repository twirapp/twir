package categories_aliases

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
)

var Remove = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "game aliase remove",
		Description: null.StringFrom("Remove category aliase"),
		RolesIDS: pq.StringArray{
			model.ChannelRoleTypeModerator.String(),
		},
		Module: "MODERATION",
		Aliases: pq.StringArray{
			"game aliase delete",
		},
		Visible: true,
		IsReply: true,
	},
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: "aliase",
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		aliase := parseCtx.ArgsParser.Get("aliase").String()

		categories, err := parseCtx.Services.CategoriesAliasesRepo.GetManyByChannelID(
			ctx,
			parseCtx.Channel.ID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get categories",
				Err:     err,
			}
		}

		var found bool

		for _, category := range categories {
			if category.Alias == aliase {
				found = true
				err = parseCtx.Services.CategoriesAliasesRepo.Delete(ctx, category.ID)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: "cannot delete category",
						Err:     err,
					}
				}
				break
			}
		}

		if !found {
			return nil, &types.CommandHandlerError{
				Message: "Category aliase not found",
			}
		}

		result := &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					"Category aliase %s removed",
					aliase,
				),
			},
		}

		return result, nil
	},
}
