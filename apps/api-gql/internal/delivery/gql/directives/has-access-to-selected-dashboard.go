package directives

import (
	"context"
	"errors"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	model "github.com/twirapp/twir/libs/gomodels"
	"gorm.io/gorm"
)

func (c *Directives) HasAccessToSelectedDashboard(
	ctx context.Context,
	obj interface{},
	next graphql.Resolver,
) (interface{}, error) {
	user, err := c.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	dashboardId, err := c.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	if user.ID == dashboardId || user.IsBotAdmin {
		return next(ctx)
	}

	var channelRoles []model.ChannelRole
	if err := c.gorm.
		WithContext(ctx).
		Where(`"channelId" = ?`, dashboardId).
		Preload("Users", `"userId" = ?`, user.ID).
		Find(&channelRoles).
		Error; err != nil {
		return nil, fmt.Errorf("cannot get channel roles: %w", err)
	}

	var userStat model.UsersStats
	if err := c.gorm.
		WithContext(ctx).
		Where(`"userId" = ? AND "channelId" = ?`, user.ID, dashboardId).
		First(&userStat).
		Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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
