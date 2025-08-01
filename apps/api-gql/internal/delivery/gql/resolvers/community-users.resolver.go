package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.76

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Masterminds/squirrel"
	helix "github.com/nicklaw5/helix/v2"
	redis "github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/dataloader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/graph"
	model "github.com/twirapp/twir/libs/gomodels"
)

// TwitchProfile is the resolver for the twitchProfile field.
func (r *communityUserResolver) TwitchProfile(ctx context.Context, obj *gqlmodel.CommunityUser) (*gqlmodel.TwirUserTwitchInfo, error) {
	return data_loader.GetHelixUserById(ctx, obj.ID)
}

// CommunityResetStats is the resolver for the communityResetStats field.
func (r *mutationResolver) CommunityResetStats(ctx context.Context, typeArg gqlmodel.CommunityUsersResetType) (bool, error) {
	user, err := r.deps.Sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	dashboardId, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	if user.ID != dashboardId {
		return false, fmt.Errorf("you cannot reset stats for this user")
	}

	var field string

	switch typeArg {
	case gqlmodel.CommunityUsersResetTypeMessages:
		field = "messages"
	case gqlmodel.CommunityUsersResetTypeWatched:
		field = "watched"
	case gqlmodel.CommunityUsersResetTypeUsedChannelsPoints:
		field = "usedChannelPoints"
	case gqlmodel.CommunityUsersResetTypeUsedEmotes:
		field = "emotes"
	}

	if field == "" {
		return false, fmt.Errorf("unknown reset typeArg: %s", typeArg)
	}

	err = r.deps.Gorm.WithContext(ctx).
		Model(&model.UsersStats{}).
		Where(`"channelId" = ?`, dashboardId).
		Update(field, 0).Error
	if err != nil {
		return false, err
	}

	iter := r.deps.Redis.Scan(
		ctx,
		0,
		fmt.Sprintf("bots:cache:ensureuser:%s:*", dashboardId),
		100,
	).Iterator()

	_, err = r.deps.Redis.Pipelined(
		ctx,
		func(pipeliner redis.Pipeliner) error {
			for iter.Next(ctx) {
				if err := pipeliner.Del(ctx, iter.Val()).Err(); err != nil {
					return err
				}
			}

			return iter.Err()
		},
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

// CommunityUsers is the resolver for the communityUsers field.
func (r *queryResolver) CommunityUsers(ctx context.Context, opts gqlmodel.CommunityUsersOpts) (*gqlmodel.CommunityUsersResponse, error) {
	var page int
	perPage := 20

	if opts.Page.IsSet() {
		page = *opts.Page.Value()
	}

	if opts.PerPage.IsSet() {
		perPage = *opts.PerPage.Value()
	}

	channel := &model.Channels{}
	err := r.deps.Gorm.
		WithContext(ctx).
		Where("channels.id = ?", opts.ChannelID).
		Joins("User").
		First(channel).Error
	if err != nil {
		return nil, err
	}

	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(
			`users_stats.*`,
		).
		From("users_stats").
		Where(
			squirrel.And{
				squirrel.Eq{`"users_stats"."channelId"`: opts.ChannelID},
				squirrel.NotEq{`"users_stats"."userId"`: opts.ChannelID},
				squirrel.NotEq{`"users_stats"."userId"`: channel.BotID},
			},
		).
		Where(`NOT EXISTS (select 1 from "users_ignored" where "id" = "users_stats"."userId")`).
		Limit(uint64(perPage)).
		Offset(uint64(page * perPage)).
		GroupBy(`"users_stats"."id"`)

	var foundTwitchChannels []helix.Channel
	if opts.Search.IsSet() {
		channels, err := r.deps.CachedTwitchClient.SearchChannels(ctx, *opts.Search.Value())
		if err != nil {
			return nil, err
		}

		foundTwitchChannels = channels
	}

	if len(foundTwitchChannels) > 0 {
		var ids []string
		for _, user := range foundTwitchChannels {
			ids = append(ids, user.ID)
		}
		queryBuilder = queryBuilder.Where(
			squirrel.Eq{
				`"users_stats"."userId"`: lo.Map(
					foundTwitchChannels,
					func(channel helix.Channel, _ int) string {
						return channel.ID
					},
				),
			},
		)
	}

	var sortBy string
	if opts.SortBy.IsSet() {
		switch *opts.SortBy.Value() {
		case gqlmodel.CommunityUsersSortByMessages:
			sortBy = "messages"
		case gqlmodel.CommunityUsersSortByUsedChannelsPoints:
			sortBy = "usedChannelPoints"
		case gqlmodel.CommunityUsersSortByUsedEmotes:
			sortBy = "emotes"
		case gqlmodel.CommunityUsersSortByWatched:
			sortBy = "watched"
		}
	}

	if sortBy != "" && !opts.Order.IsSet() {
		queryBuilder = queryBuilder.
			OrderBy(`"users_stats"."watched" DESC`)
		// GroupBy(`"users_stats"."watched"`)
	} else if sortBy != "" && opts.Order.IsSet() {
		order := *opts.Order.Value()
		queryBuilder = queryBuilder.
			OrderBy(
				fmt.Sprintf(
					`"users_stats"."%s" %s`,
					sortBy,
					strings.ToLower(order.String()),
				),
			)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("invalid query on backend: %w", err)
	}

	rows, err := r.deps.Gorm.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		r.deps.Logger.Error(
			"cannot get community users",
			slog.Any("err", err),
			slog.String("query", query),
			slog.Any("args", args),
		)
		return nil, err
	}

	var dbUsers []model.UsersStats
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
			&dbUser.CreatedAt,
			&dbUser.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		dbUsers = append(dbUsers, dbUser)
	}

	var totalStats int64
	err = r.deps.Gorm.WithContext(ctx).
		Model(&model.UsersStats{}).
		Where(`"channelId" = ?`, opts.ChannelID).
		Count(&totalStats).Error
	if err != nil {
		return nil, err
	}

	mappedUsers := make([]gqlmodel.CommunityUser, 0, len(dbUsers))
	for _, user := range dbUsers {
		mappedUsers = append(
			mappedUsers,
			gqlmodel.CommunityUser{
				ID:                user.UserID,
				WatchedMs:         int(user.Watched),
				Messages:          int(user.Messages),
				UsedEmotes:        user.Emotes,
				UsedChannelPoints: int(user.UsedChannelPoints),
			},
		)
	}

	return &gqlmodel.CommunityUsersResponse{
		Users: mappedUsers,
		Total: int(totalStats),
	}, nil
}

// CommunityUser returns graph.CommunityUserResolver implementation.
func (r *Resolver) CommunityUser() graph.CommunityUserResolver { return &communityUserResolver{r} }

type communityUserResolver struct{ *Resolver }
