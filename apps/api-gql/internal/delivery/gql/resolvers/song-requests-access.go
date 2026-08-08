package resolvers

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
)

var errSongRequestsNoChannelAccess = errors.New("user has no access to this channel song requests")

func (r *mutationResolver) ensureSongRequestsChannelAccess(ctx context.Context, channelID uuid.UUID) error {
	if channel, err := r.deps.Sessions.GetChannelFromApiKey(ctx); err == nil {
		if channel.ID != channelID {
			return errSongRequestsNoChannelAccess
		}

		return nil
	}

	user, err := r.deps.Sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return err
	}

	if r.deps.DashboardAccess == nil {
		return fmt.Errorf("dashboard access service is not configured")
	}

	hasAccess, err := r.deps.DashboardAccess.CanAccess(
		ctx,
		dashboardaccess.Subject{
			ID:         user.ID,
			IsBotAdmin: user.IsBotAdmin,
		},
		channelID,
		"MANAGE_SONG_REQUESTS",
	)
	if err != nil {
		return fmt.Errorf("cannot check dashboard access: %w", err)
	}
	if !hasAccess {
		return errSongRequestsNoChannelAccess
	}

	return nil
}
