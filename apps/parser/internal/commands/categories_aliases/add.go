package categories_aliases

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/twitch"
	categoriesaliasesrepository "github.com/twirapp/twir/libs/repositories/channels_categories_aliases"
)

var Add = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "game aliase add",
		Description: null.StringFrom("Add aliase for category for shorten usage of set command"),
		RolesIDS: pq.StringArray{
			model.ChannelRoleTypeModerator.String(),
		},
		Module: "MODERATION",
		Aliases: pq.StringArray{
			"game aliase create",
		},
		Visible: true,
		IsReply: true,
	},
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: "aliase",
		},
		command_arguments.VariadicString{
			Name: "category",
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		aliaseArg := parseCtx.ArgsParser.Get("aliase")
		categoryArg := parseCtx.ArgsParser.Get("category")

		if aliaseArg == nil || categoryArg == nil {
			return nil, &types.CommandHandlerError{
				Message: "aliase and category are required",
			}
		}

		aliase := aliaseArg.String()
		category := categoryArg.String()

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

		for _, c := range categories {
			if c.Alias == aliase {
				return nil, &types.CommandHandlerError{
					Message: fmt.Sprintf("aliase %s already exists", aliase),
				}
			}
		}

		foundTwitchCategory, err := twitch.SearchCategory(ctx, category)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get category",
				Err:     err,
			}
		}

		if foundTwitchCategory == nil {
			return nil, &types.CommandHandlerError{
				Message: "category not found",
			}
		}

		err = parseCtx.Services.CategoriesAliasesRepo.Create(
			ctx,
			categoriesaliasesrepository.CreateInput{
				ChannelID:  parseCtx.Channel.ID,
				Alias:      aliase,
				CategoryID: foundTwitchCategory.ID,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create category",
				Err:     err,
			}
		}

		result := &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					"Category aliase %s added with category %s",
					aliase,
					foundTwitchCategory.Name,
				),
			},
		}

		return result, nil
	},
}
