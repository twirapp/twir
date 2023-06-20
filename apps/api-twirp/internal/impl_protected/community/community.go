package community

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_deps"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/api/community"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Community struct {
	*impl_deps.Deps
}

func (c *Community) CommunityGetUsers(ctx context.Context, request *community.GetUsersRequest) (*community.GetUsersResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	// get sql.DB instance
	db, err := c.Db.DB()
	if err != nil {
		return nil, err
	}

	channel := &model.Channels{}
	err = c.Db.Where("id = ?", dashboardId).First(channel).Error
	if err != nil {
		return nil, err
	}

	query, args, err := squirrel.
		Select(`users_stats.*, COUNT("channels_emotes_usages"."id") AS "emotes"`).
		From("users_stats").
		LeftJoin(`"channels_emotes_usages" ON "channels_emotes_usages"."userId" = "users_stats"."userId" AND "channels_emotes_usages"."channelId" = "users_stats"."channelId"`).
		Where(squirrel.And{
			squirrel.Eq{`"users_stats"."channelId"`: dashboardId},
			// squirrel.NotEq{`"users_stats"."userId"`: channelId},
			squirrel.NotEq{`"users_stats"."userId"`: channel.BotID},
		}).
		Where(`NOT EXISTS (select 1 from "users_ignored" where "id" = "users_stats"."userId")`).
		Limit(uint64(request.Limit)).
		Offset(uint64((request.Page - 1) * request.Limit)).
		GroupBy(`"users_stats"."id"`).
		OrderBy(fmt.Sprintf(`"%s" %s`, request.SortBy.String(), request.Order.String())).
		ToSql()

	sql, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := sql.Query(args...)
	if err != nil {
		return nil, err
	}

	var dbUsers []*model.UsersStats
	for rows.Next() {
		var dbUser model.UsersStats
		err = rows.Scan(&dbUser)
		if err != nil {
			return nil, err
		}
		dbUsers = append(dbUsers, &dbUser)
	}

	return &community.GetUsersResponse{
		Users: lo.Map(dbUsers, func(item *model.UsersStats, _ int) *community.GetUsersResponse_User {
			return &community.GetUsersResponse_User{
				Id:                item.UserID,
				Watched:           item.Watched,
				Messages:          item.Messages,
				Emotes:            uint64(item.Emotes),
				UsedChannelPoints: fmt.Sprintf("%d", item.UsedChannelPoints),
			}
		}),
	}, nil
}

func (c *Community) CommunityResetStats(ctx context.Context, request *community.ResetStatsRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
