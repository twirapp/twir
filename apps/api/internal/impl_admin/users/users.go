package users

import (
	"context"
	"log/slog"
	"sync"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	messages_admin_users "github.com/twirapp/twir/libs/api/messages/admin_users"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Users struct {
	*impl_deps.Deps
}

func (c *Users) UserSwitchBan(
	ctx context.Context,
	req *messages_admin_users.UserSwitchSomeStateRequest,
) (*emptypb.Empty, error) {
	dbUser := &model.Users{}
	if err := c.Db.WithContext(ctx).Where("id = ?", req.UserId).First(dbUser).Error; err != nil {
		return nil, err
	}

	dbUser.IsBanned = !dbUser.IsBanned

	if err := c.Db.WithContext(ctx).Save(dbUser).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Users) UserSwitchAdmin(
	ctx context.Context,
	req *messages_admin_users.UserSwitchSomeStateRequest,
) (*emptypb.Empty, error) {
	dbUser := &model.Users{}
	if err := c.Db.WithContext(ctx).Where("id = ?", req.UserId).First(dbUser).Error; err != nil {
		return nil, err
	}

	dbUser.IsBotAdmin = !dbUser.IsBotAdmin

	if err := c.Db.WithContext(ctx).Save(dbUser).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Users) GetUsers(
	ctx context.Context,
	req *messages_admin_users.UsersGetRequest,
) (
	*messages_admin_users.UsersGetResponse,
	error,
) {
	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.Config, c.Grpc.Tokens)
	if err != nil {
		return nil, err
	}

	var twitchSearchUsers []helix.User
	if req.Search != nil && *req.Search != "" {
		twitchSearchUsersReq, err := twitchClient.GetUsers(&helix.UsersParams{Logins: []string{*req.Search}})
		if err != nil {
			return nil, err
		}
		if twitchSearchUsersReq.ErrorMessage != "" {
			return nil, err
		}

		twitchSearchUsers = twitchSearchUsersReq.Data.Users
	}

	page := req.GetPage()
	perPage := req.GetPerPage()
	if perPage == 0 {
		perPage = 50
	}

	var users []model.Users

	query := c.Db.
		WithContext(ctx).
		Order("id DESC").
		Joins("Channel")

	if len(twitchSearchUsers) > 0 {
		var ids []string
		for _, user := range twitchSearchUsers {
			ids = append(ids, user.ID)
		}

		query = query.Where(`"users"."id" IN ?`, ids)
	}

	if req.IsBotEnabled != nil {
		query = query.Where(`"Channel"."isEnabled" = ?`, *req.IsBotEnabled)
	}

	if req.IsAdmin != nil {
		query = query.Where(`"users"."isBotAdmin" = ?`, *req.IsAdmin)
	}

	if req.IsBanned != nil {
		query = query.Where(`"users"."is_banned" = ?`, *req.IsBanned)
	}

	var total int64
	if err := query.Model(&model.Users{}).Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.
		Limit(int(perPage)).
		Offset(int(page * perPage)).
		Find(&users).Error; err != nil {
		return nil, err
	}

	var twitchUsers []helix.User
	var twitchUsersWg sync.WaitGroup
	for _, chunk := range lo.Chunk(users, 100) {
		chunk := chunk
		twitchUsersWg.Add(1)

		go func() {
			defer twitchUsersWg.Done()

			var ids []string
			for _, user := range chunk {
				ids = append(ids, user.ID)
			}

			twitchUsersReq, err := twitchClient.GetUsers(&helix.UsersParams{IDs: ids})
			if err != nil {
				c.Logger.Error("failed to get twitch users", slog.Any("err", err))
				return
			}
			if twitchUsersReq.ErrorMessage != "" {
				c.Logger.Error("failed to get twitch users", slog.Any("err", twitchUsersReq.ErrorMessage))
				return
			}

			twitchUsers = append(twitchUsers, twitchUsersReq.Data.Users...)
		}()
	}
	twitchUsersWg.Wait()

	var mappedUsers []*messages_admin_users.UsersGetResponse_UsersGetResponseUser
	for _, user := range users {
		var twitchUser *helix.User
		for _, u := range twitchUsers {
			if u.ID == user.ID {
				twitchUser = &u
				break
			}
		}

		if twitchUser == nil {
			continue
		}

		resultedUser := &messages_admin_users.UsersGetResponse_UsersGetResponseUser{
			Id:              user.ID,
			UserName:        twitchUser.Login,
			UserDisplayName: twitchUser.DisplayName,
			Avatar:          twitchUser.ProfileImageURL,
			IsBanned:        false,
			IsAdmin:         user.IsBotAdmin,
			IsBotEnabled:    false,
		}

		if user.Channel != nil {
			resultedUser.IsBotEnabled = user.Channel.IsEnabled
			resultedUser.IsBanned = user.IsBanned
		}

		mappedUsers = append(mappedUsers, resultedUser)
	}

	return &messages_admin_users.UsersGetResponse{
		Users: mappedUsers,
		Total: int32(total),
	}, nil
}
