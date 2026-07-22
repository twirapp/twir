package directives

import (
	"context"
	"errors"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	model "github.com/twirapp/twir/libs/gomodels"
	"gorm.io/gorm"
)

func (c *Directives) HasAccessToSelectedDashboard(
	ctx context.Context,
	obj interface{},
	next graphql.Resolver,
) (interface{}, error) {
	user, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, err
	}

	dashboardId, err := c.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	legacyChannel, err := c.selectedDashboardStore.GetSelectedDashboardChannel(ctx, dashboardId)
	if err != nil {
		return nil, fmt.Errorf("cannot get channel: %w", err)
	}

	if user.IsBotAdmin {
		return next(ctx)
	}

	dashboardUUID, err := uuid.Parse(dashboardId)
	if err != nil {
		return nil, fmt.Errorf("cannot get channel: %w", err)
	}
	normalizedChannel, err := c.channels.GetChannelByID(ctx, dashboardUUID)
	if err != nil {
		return nil, fmt.Errorf("cannot get channel: %w", err)
	}

	if len(normalizedChannel.Bindings) > 0 {
		for _, binding := range normalizedChannel.Bindings {
			if binding.UserID.String() == user.ID {
				return next(ctx)
			}
		}
	} else if legacyChannel.IsOwner(user.ID) {
		return next(ctx)
	}

	channelRoles, err := c.selectedDashboardStore.GetSelectedDashboardRoles(ctx, dashboardId, user.ID)
	if err != nil {
		return nil, fmt.Errorf("cannot get channel roles: %w", err)
	}

	userStat, err := c.selectedDashboardStore.GetSelectedDashboardUserStat(ctx, user.ID, dashboardId)
	if err != nil {
		return nil, fmt.Errorf("cannot get user stats: %w", err)
	}

	roleToStats := map[model.ChannelRoleEnum]bool{
		model.ChannelRoleTypeModerator:  userStat.IsMod,
		model.ChannelRoleTypeVip:        userStat.IsVip,
		model.ChannelRoleTypeSubscriber: userStat.IsSubscriber,
	}

	for i, role := range channelRoles {
		if roleToStats[role.Type] {
			channelRoles[i].Users = append(
				role.Users,
				&model.ChannelRoleUser{
					ID:     "", // not needed
					UserID: user.ID,
					RoleID: role.ID,
				},
			)
		}
	}

	for _, role := range channelRoles {
		// we do not check does role.Users contains request author user
		// because we are doing preload by user id
		if len(role.Users) == 0 || len(role.Permissions) == 0 {
			continue
		}

		for _, roleUser := range role.Users {
			if roleUser.UserID == user.ID {
				return next(ctx)
			}
		}
	}

	return nil, fmt.Errorf("user does not have access to selected dashboard")
}

type gormSelectedDashboardStore struct {
	gorm *gorm.DB
}

func (s *gormSelectedDashboardStore) GetSelectedDashboardChannel(ctx context.Context, dashboardID string) (model.Channels, error) {
	var channel model.Channels
	if err := s.gorm.WithContext(ctx).Where("id = ?::uuid", dashboardID).First(&channel).Error; err != nil {
		return model.Channels{}, err
	}

	return channel, nil
}

func (s *gormSelectedDashboardStore) GetSelectedDashboardRoles(ctx context.Context, dashboardID, userID string) ([]model.ChannelRole, error) {
	var channelRoles []model.ChannelRole
	if err := s.gorm.
		WithContext(ctx).
		Where(`"channelId" = ?::uuid`, dashboardID).
		Preload("Users", `user_id = ?`, userID).
		Find(&channelRoles).
		Error; err != nil {
		return nil, err
	}

	return channelRoles, nil
}

func (s *gormSelectedDashboardStore) GetSelectedDashboardUserStat(ctx context.Context, userID, dashboardID string) (model.UsersStats, error) {
	var userStat model.UsersStats
	err := s.gorm.
		WithContext(ctx).
		Where(`user_id = ? AND channel_id = ?::uuid`, userID, dashboardID).
		First(&userStat).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return userStat, nil
	}
	if err != nil {
		return model.UsersStats{}, err
	}

	return userStat, nil
}
