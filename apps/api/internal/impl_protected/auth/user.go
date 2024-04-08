package auth

import (
	"context"
	"fmt"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/helpers"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/auth"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Auth struct {
	*impl_deps.Deps
}

func (c *Auth) AuthUserProfile(ctx context.Context, _ *emptypb.Empty) (*auth.Profile, error) {
	dbUser, err := helpers.GetUserModelFromCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get user model from ctx: %w", err)
	}

	twitchUser := c.SessionManager.Get(ctx, "twitchUser").(helix.User)
	selectedDashboardId := c.SessionManager.Get(ctx, "dashboardId").(string)

	return &auth.Profile{
		Id:                  dbUser.ID,
		Avatar:              twitchUser.ProfileImageURL,
		Login:               twitchUser.Login,
		DisplayName:         twitchUser.DisplayName,
		ApiKey:              dbUser.ApiKey,
		IsBotAdmin:          dbUser.IsBotAdmin,
		SelectedDashboardId: selectedDashboardId,
		HideOnLandingPage:   dbUser.HideOnLandingPage,
	}, nil
}

func (c *Auth) AuthSetDashboard(ctx context.Context, req *auth.SetDashboard) (
	*emptypb.Empty,
	error,
) {
	dbUser := c.SessionManager.Get(ctx, "dbUser").(model.Users)

	var usersRoles []*model.ChannelRoleUser
	if err := c.Db.Where(
		`"userId" = ?`,
		dbUser.ID,
	).Preload("Role").Find(&usersRoles).Error; err != nil {
		return nil, err
	}

	var channelRoles []*model.ChannelRole
	if err := c.Db.Where(`"channelId" = ?`, req.DashboardId).Find(&channelRoles).Error; err != nil {
		return nil, err
	}

	usersStats := model.UsersStats{}
	if err := c.Db.Where(
		`"userId" = ? AND "channelId" = ?`,
		dbUser.ID,
		req.DashboardId,
	).Find(&usersStats).Error; err != nil {
		return nil, err
	}

	hasPermission := lo.SomeBy(
		usersRoles, func(role *model.ChannelRoleUser) bool {
			return role.UserID == dbUser.ID && role.Role.ChannelID == req.DashboardId
		},
	) || lo.SomeBy(
		channelRoles, func(role *model.ChannelRole) bool {
			if role.Type == model.ChannelRoleTypeModerator {
				return usersStats.IsMod
			} else if role.Type == model.ChannelRoleTypeVip {
				return usersStats.IsVip
			} else if role.Type == model.ChannelRoleTypeSubscriber {
				return usersStats.IsSubscriber
			} else {
				return false
			}
		},
	)

	if !hasPermission && !dbUser.IsBotAdmin && dbUser.ID != req.DashboardId {
		return nil, fmt.Errorf(
			"user %s does not have permission to access dashboard %s",
			dbUser.ID,
			req.DashboardId,
		)
	}

	c.SessionManager.Put(ctx, "dashboardId", req.DashboardId)

	return &emptypb.Empty{}, nil
}

func (c *Auth) AuthGetDashboards(
	ctx context.Context,
	_ *emptypb.Empty,
) (*auth.GetDashboardsResponse, error) {
	dbUser := c.SessionManager.Get(ctx, "dbUser").(model.Users)
	var dashboards []*auth.Dashboard

	if dbUser.IsBotAdmin {
		var channels []*model.Channels
		if err := c.Db.Find(&channels).Error; err != nil {
			return nil, err
		}

		for _, channel := range channels {
			dashboards = append(
				dashboards, &auth.Dashboard{
					Id:    channel.ID,
					Flags: []string{model.RolePermissionCanAccessDashboard.String()},
				},
			)
		}
	} else {
		var roles []*model.ChannelRoleUser
		if err := c.Db.Where(`"userId" = ?`, dbUser.ID).Preload("Role").Find(&roles).Error; err != nil {
			return nil, err
		}
		for _, role := range roles {
			dashboards = append(
				dashboards, &auth.Dashboard{
					Id:    role.Role.ChannelID,
					Flags: role.Role.Permissions,
				},
			)
		}
	}

	var usersStats []model.UsersStats
	if err := c.Db.Where(`"userId" = ?`, dbUser.ID).Find(&usersStats).Error; err != nil {
		return nil, err
	}

	for _, i := range usersStats {
		var channelRoles []model.ChannelRole
		if err := c.Db.Where(`"channelId" = ?`, i.ChannelID).Find(&channelRoles).Error; err != nil {
			return nil, err
		}

		var role model.ChannelRole

		if i.IsMod {
			role, _ = lo.Find(
				channelRoles, func(role model.ChannelRole) bool {
					return role.Type == model.ChannelRoleTypeModerator
				},
			)
		} else if i.IsVip {
			role, _ = lo.Find(
				channelRoles, func(role model.ChannelRole) bool {
					return role.Type == model.ChannelRoleTypeVip
				},
			)
		} else if i.IsSubscriber {
			role, _ = lo.Find(
				channelRoles, func(role model.ChannelRole) bool {
					return role.Type == model.ChannelRoleTypeSubscriber
				},
			)
		}

		if role.ID != "" {
			dashboards = append(
				dashboards, &auth.Dashboard{
					Id:    role.ChannelID,
					Flags: role.Permissions,
				},
			)
		}
	}

	dashboards = lo.UniqBy(
		dashboards,
		func(dashboard *auth.Dashboard) string {
			return dashboard.Id
		},
	)

	dashboards = lo.Filter(
		dashboards,
		func(dashboard *auth.Dashboard, _ int) bool {
			return len(dashboard.Flags) > 0
		},
	)

	return &auth.GetDashboardsResponse{
		Dashboards: dashboards,
	}, nil
}

func (c *Auth) AuthLogout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	c.SessionManager.Destroy(ctx)

	return &emptypb.Empty{}, nil
}
