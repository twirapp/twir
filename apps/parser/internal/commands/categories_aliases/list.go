package categories_aliases

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/twitch"
)

var List = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "game aliase list",
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
			parseCtx.Channel.ID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get categories",
				Err:     err,
			}
		}

		if len(categories) == 0 {
			return &types.CommandsHandlerResult{
				Result: []string{"No categories aliases created."},
			}, nil
		}

		twitchClient, err := twitch.NewAppClientWithContext(
			ctx,
			*parseCtx.Services.Config,
			parseCtx.Services.Bus,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create twitch client",
				Err:     err,
			}
		}

		categoriesIds := make([]string, len(categories))
		for i, category := range categories {
			categoriesIds[i] = category.CategoryID
		}

		gamesRequest, err := twitchClient.GetGames(
			&helix.GamesParams{
				IDs: categoriesIds,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get games",
				Err:     err,
			}
		}
		if gamesRequest.ErrorMessage != "" {
			return nil, &types.CommandHandlerError{
				Message: "cannot get games",
				Err:     fmt.Errorf(gamesRequest.ErrorMessage),
			}
		}

		aliases := make([]createdAliase, 0, len(categories))
		for idx, category := range categories {
			aliases = append(
				aliases, createdAliase{
					aliase: category.Alias,
				},
			)

			for _, game := range gamesRequest.Data.Games {
				if game.ID == category.CategoryID {
					aliases[idx].twitchCategory = &game
					break
				}
			}
		}

		slices.SortFunc(
			aliases, func(a, b createdAliase) int {
				return strings.Compare(a.aliase, b.aliase)
			},
		)

		var resultedString strings.Builder

		for _, aliase := range aliases {
			if resultedString.Len() > 0 {
				resultedString.WriteString(" · ")
			}

			resultedString.WriteString(aliase.aliase)

			if aliase.twitchCategory != nil {
				resultedString.WriteString(fmt.Sprintf(" (%s)", aliase.twitchCategory.Name))
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
	aliase         string
	twitchCategory *helix.Game
}
