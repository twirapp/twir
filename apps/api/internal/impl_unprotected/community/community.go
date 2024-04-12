package community

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/community"
	"github.com/twitchtv/twirp"
)

type Community struct {
	*impl_deps.Deps
}

func (c *Community) CommunityGetUsers(
	ctx context.Context,
	request *community.GetUsersRequest,
) (*community.GetUsersResponse, error) {
	if request.ChannelId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "channelId is required")
	}

	// get sql.DB instance
	db, err := c.Db.DB()
	if err != nil {
		return nil, err
	}

	channel := &model.Channels{}
	err = c.Db.WithContext(ctx).Where("channels.id = ?", request.ChannelId).Joins("User").First(
		channel,
	).Error
	if err != nil {
		return nil, err
	}

	if channel.User.IsBanned {
		return &community.GetUsersResponse{}, nil
	}

	var sortBy string
	if request.SortBy == community.GetUsersRequest_UsedChannelPoints {
		sortBy = "usedChannelPoints"
	} else {
		sortBy = strings.ToLower(request.SortBy.String())
	}

	// orderBy := fmt.Sprintf(`"users_stats"."%s"`, sortBy)
	// if request.SortBy == community.GetUsersRequest_Emotes {
	//	orderBy = "emotes"
	// }

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(`users_stats.*, COUNT("channels_emotes_usages"."id") AS "emotes"`).
		From("users_stats").
		LeftJoin(`"channels_emotes_usages" ON "channels_emotes_usages"."userId" = "users_stats"."userId" AND "channels_emotes_usages"."channelId" = "users_stats"."channelId"`).
		Where(
			squirrel.And{
				squirrel.Eq{`"users_stats"."channelId"`: request.ChannelId},
				squirrel.NotEq{`"users_stats"."userId"`: request.ChannelId},
				squirrel.NotEq{`"users_stats"."userId"`: channel.BotID},
				squirrel.Gt{`"users_stats"."messages"`: 0},
			},
		).
		Where(`NOT EXISTS (select 1 from "users_ignored" where "id" = "users_stats"."userId")`).
		Limit(uint64(request.Limit)).
		Offset(uint64((request.Page - 1) * request.Limit)).
		GroupBy(`"users_stats"."id"`).
		OrderBy(fmt.Sprintf(`"%s" %s`, sortBy, strings.ToLower(request.Order.String()))).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var dbUsers []*model.UsersStats
	for rows.Next() {
		var dbUser model.UsersStats

		err = rows.Scan(
			&dbUser.ID,
			&dbUser.Messages,
			&dbUser.Watched,
			&dbUser.ChannelID,
			&dbUser.UserID,
			&dbUser.UsedChannelPoints,
			&dbUser.IsMod,
			&dbUser.IsVip,
			&dbUser.IsSubscriber,
			&dbUser.Reputation,
			&dbUser.Emotes,
		)
		if err != nil {
			return nil, err
		}
		dbUsers = append(dbUsers, &dbUser)
	}

	var totalStats int64
	err = c.Db.WithContext(ctx).
		Model(&model.UsersStats{}).
		Where(`"channelId" = ? AND "messages" > 0`, request.ChannelId).
		Count(&totalStats).Error
	if err != nil {
		return nil, err
	}

	// totalPages := (totalStats + int64(request.Limit) - 1) / int64(request.Limit)

	return &community.GetUsersResponse{
		Users: lo.Map(
			dbUsers, func(item *model.UsersStats, _ int) *community.GetUsersResponse_User {
				return &community.GetUsersResponse_User{
					Id:                item.UserID,
					Watched:           fmt.Sprint(item.Watched),
					Messages:          item.Messages,
					Emotes:            fmt.Sprint(item.Emotes),
					UsedChannelPoints: fmt.Sprint(item.UsedChannelPoints),
				}
			},
		),
		TotalUsers: uint32(totalStats),
	}, nil
}
