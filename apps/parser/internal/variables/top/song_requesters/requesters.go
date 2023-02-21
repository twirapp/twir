package top_song_requesters

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	config "github.com/satont/tsuwari/libs/config"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

type Request struct {
	ChannelID   string
	OrderedById string
	Count       int
}

var Variable = types.Variable{
	Name:        "top.songs.requesters",
	Description: lo.ToPtr("Top users by requested songs"),
	Example:     lo.ToPtr("top.songs.requesters|10"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		sqlxDb := do.MustInvoke[sqlx.DB](di.Provider)
		cfg := do.MustInvoke[config.Config](di.Provider)
		tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
		twitchClient, err := twitch.NewAppClient(cfg, tokensGrpc)

		result := &types.VariableHandlerResult{}
		var page = 1

		if ctx.Text != nil {
			p, err := strconv.Atoi(*ctx.Text)
			if err != nil {
				page = p
			}

			if page <= 0 {
				page = 1
			}
		}

		limit := 10
		if data.Params != nil {
			newLimit, err := strconv.Atoi(*data.Params)
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
			Where(`"channelId" = ?`, ctx.ChannelId).
			GroupBy(`"orderedById"`).
			OrderBy("count desc").
			Limit(uint64(limit)).
			ToSql()
		query = sqlxDb.Rebind(query)
		fmt.Println(query)
		err = sqlxDb.Select(&dbEntities, query, args...)

		if err != nil {
			fmt.Println(err)
			result.Result = "internal error"
			return result, nil
		}

		twitchUsers, err := twitchClient.GetUsers(&helix.UsersParams{
			IDs: lo.Map(dbEntities, func(item model.RequestedSong, _ int) string {
				return item.OrderedById
			}),
		})

		if err != nil {
			return nil, err
		}

		mappedTop := []string{}

		for _, entity := range dbEntities {
			user, ok := lo.Find(twitchUsers.Data.Users, func(item helix.User) bool {
				return item.ID == entity.OrderedById
			})

			if !ok {
				continue
			}

			mappedTop = append(mappedTop, fmt.Sprintf(
				"%s × %v",
				user.Login,
				entity.Count,
			))
		}

		result.Result = strings.Join(mappedTop, " · ")
		return result, nil
	},
}
