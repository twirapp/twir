package auth

import (
	"context"

	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/auth"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Auth struct {
	*impl_deps.Deps
}

func (c *Auth) AuthUserProfile(ctx context.Context, _ *emptypb.Empty) (*auth.Profile, error) {
	dbUser := c.SessionManager.Get(ctx, "dbUser").(model.Users)
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
	}, nil
}

func (c *Auth) AuthSetDashboard(ctx context.Context, req *auth.SetDashboard) (*emptypb.Empty, error) {
	c.SessionManager.Put(ctx, "dashboardId", req.DashboardId)

	return &emptypb.Empty{}, nil
}

func (c *Auth) AuthGetDashboards(ctx context.Context, _ *emptypb.Empty) (*auth.GetDashboardsResponse, error) {
	dbUser := c.SessionManager.Get(ctx, "dbUser").(model.Users)
	var dashboards []*auth.Dashboard

	if dbUser.IsBotAdmin {
		var channels []*model.Channels
		if err := c.Db.Find(&channels).Error; err != nil {
			return nil, err
		}

		for _, channel := range channels {
			dashboards = append(dashboards, &auth.Dashboard{
				Id:    channel.ID,
				Flags: []string{model.RolePermissionCanAccessDashboard.String()},
			})
		}
	} else {
		var roles []*model.ChannelRoleUser
		if err := c.Db.Where(`"userId" = ?`, dbUser.ID).Preload("Role").Find(&roles).Error; err != nil {
			return nil, err
		}
		for _, role := range roles {
			dashboards = append(dashboards, &auth.Dashboard{
				Id:    role.Role.ChannelID,
				Flags: role.Role.Permissions,
			})
		}
	}

	return &auth.GetDashboardsResponse{
		Dashboards: dashboards,
	}, nil
}

func (c *Auth) AuthLogout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	c.SessionManager.Destroy(ctx)

	return &emptypb.Empty{}, nil
}
