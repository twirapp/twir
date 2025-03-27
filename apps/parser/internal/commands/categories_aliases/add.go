package categories_aliases

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
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

		twitchClient, err := twitch.NewAppClientWithContext(
			ctx,
			*parseCtx.Services.Config,
			parseCtx.Services.GrpcClients.Tokens,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create twitch client",
				Err:     err,
			}
		}

		categoriesResponse, err := twitchClient.SearchCategories(
			&helix.SearchCategoriesParams{
				Query: category,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get categories",
				Err:     err,
			}
		}
		if categoriesResponse.ErrorMessage != "" {
			return nil, &types.CommandHandlerError{
				Message: "cannot get categories",
				Err:     fmt.Errorf(categoriesResponse.ErrorMessage),
			}
		}

		if len(categoriesResponse.Data.Categories) == 0 {
			return nil, &types.CommandHandlerError{
				Message: "category not found",
			}
		}

		twitchCategory := categoriesResponse.Data.Categories[0]

		err = parseCtx.Services.CategoriesAliasesRepo.Create(
			ctx,
			categoriesaliasesrepository.CreateInput{
				ChannelID:  parseCtx.Channel.ID,
				Alias:      aliase,
				CategoryID: twitchCategory.ID,
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
					twitchCategory.Name,
				),
			},
		}

		return result, nil
	},
}
