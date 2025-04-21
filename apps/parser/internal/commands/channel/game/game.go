package channel_game

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"

	"github.com/nicklaw5/helix/v2"
)

const (
	gameArgName = "gameOrAliase"
)

var SetCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "game",
		Description: null.StringFrom("Change category of channel"),
		Module:      "MODERATION",
		IsReply:     true,
		Visible:     false,
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
	},
	Args: []command_arguments.Arg{
		command_arguments.VariadicString{
			Name: gameArgName,
			Hint: "category name or created category aliase",
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		categoryArg := parseCtx.ArgsParser.Get(gameArgName).String()

		categoryAliases, err := parseCtx.Services.CategoriesAliasesRepo.GetManyByChannelID(
			ctx,
			parseCtx.Channel.ID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get category aliases",
				Err:     err,
			}
		}

		twitchClient, err := twitch.NewUserClientWithContext(
			ctx,
			parseCtx.Channel.ID,
			*parseCtx.Services.Config,
			parseCtx.Services.GrpcClients.Tokens,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create broadcaster twitch client",
				Err:     err,
			}
		}

		for _, categoryAlias := range categoryAliases {
			if categoryAlias.Alias == categoryArg {
				changeResponse, err := twitchClient.EditChannelInformation(
					&helix.EditChannelInformationParams{
						BroadcasterID: parseCtx.Channel.ID,
						GameID:        categoryAlias.CategoryID,
					},
				)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: "cannot change category",
						Err:     err,
					}
				}
				if changeResponse.ErrorMessage != "" {
					return nil, &types.CommandHandlerError{
						Message: fmt.Sprintf("cannot change category: %s", changeResponse.ErrorMessage),
						Err:     fmt.Errorf(changeResponse.ErrorMessage),
					}
				}

				categoryRequest, err := twitchClient.GetGames(
					&helix.GamesParams{
						IDs: []string{categoryAlias.CategoryID},
					},
				)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: "cannot get category",
						Err:     err,
					}
				}
				if categoryRequest.ErrorMessage != "" {
					return nil, &types.CommandHandlerError{
						Message: fmt.Sprintf("cannot get category: %s", categoryRequest.ErrorMessage),
						Err:     fmt.Errorf(categoryRequest.ErrorMessage),
					}
				}

				if len(categoryRequest.Data.Games) == 0 {
					return nil, &types.CommandHandlerError{
						Message: "category not found",
						Err:     fmt.Errorf("category not found"),
					}
				}

				result.Result = append(
					result.Result,
					fmt.Sprintf("✅ %s", categoryRequest.Data.Games[0].Name),
				)
				return result, nil
			}
		}

		category, err := parseCtx.Services.CacheTwitchClient.SearchCategory(ctx, categoryArg)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "game not found on twitch",
				Err:     err,
			}
		}

		changeResponse, err := twitchClient.EditChannelInformation(
			&helix.EditChannelInformationParams{
				BroadcasterID: parseCtx.Channel.ID,
				GameID:        category.ID,
			},
		)

		if err != nil || changeResponse.StatusCode != 204 {
			result.Result = append(
				result.Result,
				lo.If(changeResponse.ErrorMessage != "", changeResponse.ErrorMessage).Else(
					"❌ internal error",
				),
			)
			return result, nil
		}

		result.Result = append(result.Result, "✅ "+category.Name)
		return result, nil
	},
}
