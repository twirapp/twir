package channel_game

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"github.com/twirapp/twir/libs/twitch"

	"github.com/nicklaw5/helix/v2"
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
			Name:     gameArgName,
			Hint:     i18n.Get(locales.Translations.Commands.Channel.Hints.GameArgName),
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
			parseCtx.Channel.ID,
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
			channelInfo, err := twitchClient.GetChannelInformation(
				&helix.GetChannelInformationParams{
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
			if len(channelInfo.Data.Channels) == 0 {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(
						ctx,
						locales.Translations.Commands.Channel.Errors.ChannelNotFound,
					),
					Err: fmt.Errorf(i18n.GetCtx(ctx, locales.Translations.Commands.Channel.Errors.ChannelNotFound)),
				}
			}

			result.Result = append(result.Result, channelInfo.Data.Channels[0].GameName)
			return result, nil
		}

		categoryArg := parseCtx.ArgsParser.Get(gameArgName).String()

		categoryAliases, err := parseCtx.Services.CategoriesAliasesRepo.GetManyByChannelID(
			ctx,
			parseCtx.Channel.ID,
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
				changeResponse, err := twitchClient.EditChannelInformation(
					&helix.EditChannelInformationParams{
						BroadcasterID: parseCtx.Channel.ID,
						GameID:        categoryAlias.CategoryID,
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
				if changeResponse.ErrorMessage != "" {
					return nil, &types.CommandHandlerError{
						Message: i18n.GetCtx(
							ctx,
							locales.Translations.Commands.Channel.Errors.CategoryCannotChangeError.
								SetVars(locales.KeysCommandsChannelErrorsCategoryCannotChangeErrorVars{ErrorMessage: changeResponse.ErrorMessage}),
						),
						Err: fmt.Errorf(changeResponse.ErrorMessage),
					}
				}

				categoryRequest, err := twitchClient.GetGames(
					&helix.GamesParams{
						IDs: []string{categoryAlias.CategoryID},
					},
				)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: i18n.GetCtx(ctx, locales.Translations.Commands.Channel.Errors.CategoryCannotGet),
						Err:     err,
					}
				}
				if categoryRequest.ErrorMessage != "" {
					return nil, &types.CommandHandlerError{
						Message: i18n.GetCtx(
							ctx,
							locales.Translations.Commands.Channel.Errors.CategoryCannotGetError.
								SetVars(locales.KeysCommandsChannelErrorsCategoryCannotGetErrorVars{ErrorMessage: categoryRequest.ErrorMessage}),
						),
						Err: fmt.Errorf(categoryRequest.ErrorMessage),
					}
				}

				if len(categoryRequest.Data.Games) == 0 {
					return nil, &types.CommandHandlerError{
						Message: i18n.GetCtx(
							ctx,
							locales.Translations.Commands.Channel.Errors.CategoryNotFound,
						),
						Err: fmt.Errorf(i18n.GetCtx(
							ctx,
							locales.Translations.Commands.Channel.Errors.CategoryNotFound,
						)),
					}
				}

				result.Result = append(
					result.Result,
					i18n.GetCtx(
						ctx,
						locales.Translations.Commands.Channel.Add.CategoryChange.
							SetVars(locales.KeysCommandsChannelAddCategoryChangeVars{CategoryName: categoryRequest.Data.Games[0].Name}),
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
					i18n.GetCtx(
						ctx,
						locales.Translations.Commands.Channel.Errors.Internal,
					),
				),
			)
			return result, nil
		}

		result.Result = append(result.Result, "✅ "+category.Name)
		return result, nil
	},
}
