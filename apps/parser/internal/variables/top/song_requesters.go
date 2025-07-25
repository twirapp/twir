package top

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/twitch"
	"go.uber.org/zap"
)

var SongRequesters = &types.Variable{
	Name:                "top.songs.requesters",
	Description:         lo.ToPtr("Top users by requested songs"),
	Example:             lo.ToPtr("top.songs.requesters|10"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		twitchClient, err := twitch.NewAppClientWithContext(
			ctx,
			*parseCtx.Services.Config,
			parseCtx.Services.Bus,
		)

		result := &types.VariableHandlerResult{}
		page := 1

		if parseCtx.Text != nil {
			p, err := strconv.Atoi(*parseCtx.Text)
			if err != nil {
				page = p
			}

			if page <= 0 {
				page = 1
			}
		}

		limit := 10
		if variableData.Params != nil {
			newLimit, err := strconv.Atoi(*variableData.Params)
			if err == nil {
				limit = newLimit
			}
		}

		if limit > 50 {
			limit = 10
		}

		dbEntities := []model.RequestedSong{}

		query, args, err := squirrel.
			Select(`"orderedById"`, "COUNT(*) as count").
			From("channels_requested_songs").
			Where(`"channelId" = ?`, parseCtx.Channel.ID).
			GroupBy(`"orderedById"`).
			OrderBy("count desc").
			Limit(uint64(limit)).
			ToSql()
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			return result, nil
		}

		query = parseCtx.Services.Sqlx.Rebind(query)
		err = parseCtx.Services.Sqlx.Select(&dbEntities, query, args...)

		if err != nil {
			zap.S().Error(err)
			result.Result = "internal error"
			return result, nil
		}

		twitchUsers, err := twitchClient.GetUsers(
			&helix.UsersParams{
				IDs: lo.Map(
					dbEntities, func(item model.RequestedSong, _ int) string {
						return item.OrderedById
					},
				),
			},
		)
		if err != nil {
			return nil, err
		}

		mappedTop := []string{}

		for _, entity := range dbEntities {
			user, ok := lo.Find(
				twitchUsers.Data.Users, func(item helix.User) bool {
					return item.ID == entity.OrderedById
				},
			)

			if !ok {
				continue
			}

			mappedTop = append(
				mappedTop, fmt.Sprintf(
					"%s × %v",
					user.Login,
					entity.Count,
				),
			)
		}

		result.Result = strings.Join(mappedTop, " · ")
		return result, nil
	},
}
