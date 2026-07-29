package categories_aliases

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/guregu/null"
	"github.com/kvizyx/twitchy/helix"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"github.com/twirapp/twir/libs/twitch"
)

var List = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "game alias list",
		Description: null.StringFrom("List created categories aliases"),
		RolesIDS: pq.StringArray{
			model.ChannelRoleTypeModerator.String(),
		},
		Module:  "MODERATION",
		Aliases: pq.StringArray{},
		Visible: true,
		IsReply: true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		categories, err := parseCtx.Services.CategoriesAliasesRepo.GetManyByChannelID(
			ctx,
			parseCtx.Channel.DBChannelID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.CategoriesAliases.Errors.CategoryCannotToGet,
				),
				Err: err,
			}
		}

		if len(categories) == 0 {
			return &types.CommandsHandlerResult{
				Result: []string{i18n.GetCtx(
					ctx,
					locales.Translations.Commands.CategoriesAliases.Errors.AliasEmpty,
				)},
			}, nil
		}

		twitchClient, err := twitch.NewAppClientWithContext(
			ctx,
			*parseCtx.Services.Config,
			parseCtx.Services.Bus,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.CategoriesAliases.Errors.TwitchClientCannotToCreate,
				),
				Err: err,
			}
		}

		categoriesIds := make([]string, len(categories))
		for i, category := range categories {
			categoriesIds[i] = category.CategoryID
		}

		gamesRequest, err := twitchClient.Games.GetGames(
			ctx,
			helix.GetGamesRequest{
				IDs: categoriesIds,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.CategoriesAliases.Errors.GameCannotToGet,
				),
				Err: err,
			}
		}
		aliases := make([]createdAliase, 0, len(categories))
		for idx, category := range categories {
			aliases = append(
				aliases, createdAliase{
					alias: category.Alias,
				},
			)

			for _, game := range gamesRequest.Data {
				if game.ID == category.CategoryID {
					aliases[idx].twitchCategory = &game
					break
				}
			}
		}

		slices.SortFunc(
			aliases, func(a, b createdAliase) int {
				return strings.Compare(a.alias, b.alias)
			},
		)

		var resultedString strings.Builder

		for _, alias := range aliases {
			if resultedString.Len() > 0 {
				resultedString.WriteString(" · ")
			}

			resultedString.WriteString(alias.alias)

			if alias.twitchCategory != nil {
				resultedString.WriteString(fmt.Sprintf(" (%s)", alias.twitchCategory.Name))
			} else {
				resultedString.WriteString(" (not found)")
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{resultedString.String()},
		}, nil
	},
}

type createdAliase struct {
	alias          string
	twitchCategory *helix.Game
}
