package resolvers

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var errCannotUnlinkCurrentPlatform = errors.New("cannot unlink current platform")

func (r *Resolver) selectedChannelPlatformDashboard(ctx context.Context) (uuid.UUID, error) {
	if r.deps.ChannelPlatformDashboard == nil {
		return uuid.Nil, fmt.Errorf("channel platform dashboard session is not configured")
	}

	dashboardID, err := r.deps.ChannelPlatformDashboard.GetSelectedDashboard(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("get selected dashboard: %w", err)
	}
	parsedDashboardID, err := uuid.Parse(dashboardID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse selected dashboard: %w", err)
	}

	return parsedDashboardID, nil
}
