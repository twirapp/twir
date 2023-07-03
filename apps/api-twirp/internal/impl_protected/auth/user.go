package auth

import (
	"context"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/api-twirp/internal/impl_deps"
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
