package top

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	config "github.com/satont/tsuwari/libs/config"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"

	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
)

type UserStats struct {
	DisplayName string
	UserName    string
	Value       int
}

func GetTop(
	ctx *variables_cache.VariablesCacheService,
	channelId string,
	topType string,
	page *int,
) []*UserStats {
	cfg := do.MustInvoke[config.Config](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	sqlxDb := do.MustInvoke[sqlx.DB](di.Provider)

	twitchClient, err := twitch.NewAppClient(cfg, tokensGrpc)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	if page == nil {
		newPage := 1
		page = &newPage
	}

	offset := (*page - 1) * 10

	// another approach how to filter is via joins, but i decided to leave it with sub query
	//LEFT JOIN "users_ignored" ON "users_ignored"."id" = "users_stats"."userId"
	//WHERE
	//"users_stats"."channelId" = '869882828' AND "users_ignored"."id" is null
	query, args, err := squirrel.
		Select("users_stats.*").
		From("users_stats").
		Where(squirrel.Eq{`"channelId"`: channelId}).
		Where(`NOT EXISTS (select 1 from users_ignored where "id" = "users_stats"."userId")`).
		Limit(10).
		Offset(uint64(offset)).
		OrderBy(fmt.Sprintf("%s DESC", topType)).
		ToSql()
	query = sqlxDb.Rebind(query)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	records := []model.UsersStats{}

	err = sqlxDb.Select(&records, query, args...)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	ids := lo.Map(records, func(record model.UsersStats, _ int) string {
		return record.UserID
	})

	twitchUsers, err := twitchClient.GetUsers(&helix.UsersParams{
		IDs: ids,
	})

	if err != nil || len(twitchUsers.Data.Users) == 0 {
		return nil
	}

	stats := lo.Map(records, func(record model.UsersStats, _ int) *UserStats {
		twitchUser, ok := lo.Find(twitchUsers.Data.Users, func(user helix.User) bool {
			return user.ID == record.UserID
		})

		if !ok {
			return nil
		}

		res := &UserStats{
			DisplayName: twitchUser.DisplayName,
			UserName:    twitchUser.Login,
		}

		if topType == "messages" {
			res.Value = int(record.Messages)
		}

		if topType == "watched" {
			res.Value = int(record.Watched)
		}

		return res
	})

	return stats
}
