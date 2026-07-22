package resolvers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

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
