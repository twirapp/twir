package channel_game

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/kvizyx/twitchy/helix"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"github.com/twirapp/twir/libs/twitch"
)

const (
	gameArgName = "gameOrAlias"
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
			HintFunc: func(ctx context.Context) string {
				return i18n.GetCtx(ctx, locales.Translations.Commands.Channel.Hints.GameArgName)
			},
			Optional: true,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		twitchClient, err := twitch.NewUserClientWithContext(
			ctx,
			parseCtx.Channel.TwitchUserID,
			*parseCtx.Services.Config,
			parseCtx.Services.Bus,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Channel.Errors.BroadcasterTwitchClientCannotCreate,
				),
				Err: err,
			}
		}

		if !parseCtx.ArgsParser.IsExists(gameArgName) {
			channelInfo, err := twitchClient.Channels.GetChannelInformation(
				ctx,
				helix.GetChannelInformationRequest{
					BroadcasterIDs: []string{parseCtx.Channel.ID},
				},
			)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(
						ctx,
						locales.Translations.Commands.Channel.Errors.ChannelCannotGetInformation,
					),
					Err: err,
				}
			}
			if len(channelInfo.Data) == 0 {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(
						ctx,
						locales.Translations.Commands.Channel.Errors.ChannelNotFound,
					),
					Err: fmt.Errorf(
						i18n.GetCtx(
							ctx,
							locales.Translations.Commands.Channel.Errors.ChannelNotFound,
						),
					),
				}
			}

			result.Result = append(result.Result, channelInfo.Data[0].GameName)
			return result, nil
		}

		categoryArg := parseCtx.ArgsParser.Get(gameArgName).String()

		categoryAliases, err := parseCtx.Services.CategoriesAliasesRepo.GetManyByChannelID(
			ctx,
			parseCtx.Channel.DBChannelID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Channel.Errors.AliasCannotGetCategory,
				),
				Err: err,
			}
		}

		for _, categoryAlias := range categoryAliases {
			if categoryAlias.Alias == categoryArg {
				gameID := categoryAlias.CategoryID
				_, err := twitchClient.Channels.ModifyChannelInformation(
					ctx,
					helix.ModifyChannelInformationRequest{
						BroadcasterID: parseCtx.Channel.ID,
						GameID:        &gameID,
					},
				)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: i18n.GetCtx(
							ctx,
							locales.Translations.Commands.Channel.Errors.CategoryCannotChange,
						),
						Err: err,
					}
				}
				categoryRequest, err := twitchClient.Games.GetGames(
					ctx,
					helix.GetGamesRequest{
						IDs: []string{categoryAlias.CategoryID},
					},
				)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: i18n.GetCtx(
							ctx,
							locales.Translations.Commands.Channel.Errors.CategoryCannotGet,
						),
						Err: err,
					}
				}
				if len(categoryRequest.Data) == 0 {
					return nil, &types.CommandHandlerError{
						Message: i18n.GetCtx(
							ctx,
							locales.Translations.Commands.Channel.Errors.CategoryNotFound,
						),
						Err: fmt.Errorf(
							i18n.GetCtx(
								ctx,
								locales.Translations.Commands.Channel.Errors.CategoryNotFound,
							),
						),
					}
				}

				result.Result = append(
					result.Result,
					i18n.GetCtx(
						ctx,
						locales.Translations.Commands.Channel.Add.CategoryChange.
							SetVars(locales.KeysCommandsChannelAddCategoryChangeVars{CategoryName: categoryRequest.Data[0].Name}),
					),
				)
				return result, nil
			}
		}

		category, err := parseCtx.Services.CacheTwitchClient.SearchCategory(ctx, categoryArg)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Channel.Errors.GameNotFound,
				),
				Err: err,
			}
		}

		categoryID := category.ID
		_, err = twitchClient.Channels.ModifyChannelInformation(
			ctx,
			helix.ModifyChannelInformationRequest{
				BroadcasterID: parseCtx.Channel.ID,
				GameID:        &categoryID,
			},
		)

		if err != nil {
			result.Result = append(
				result.Result,
				i18n.GetCtx(ctx, locales.Translations.Errors.Generic.Internal),
			)
			return result, nil
		}

		result.Result = append(result.Result, "✅ "+category.Name)
		return result, nil
	},
}
