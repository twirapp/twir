package resolvers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlerrors"
	dotaentity "github.com/twirapp/twir/libs/entities/dota"
)

func getDotaDashboardID(ctx context.Context, sessions SessionReader) (uuid.UUID, error) {
	dashboardID, err := sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return uuid.Nil, gqlerrors.HandleError(err)
	}

	parsed, err := uuid.Parse(dashboardID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid selected dashboard id: %w", err)
	}

	return parsed, nil
}

func getDotaSettings(ctx context.Context, deps *Deps) (dotaentity.ChannelDotaSettings, error) {
	dashboardID, err := getDotaDashboardID(ctx, deps.Sessions)
	if err != nil {
		return dotaentity.Nil, err
	}

	settings, err := deps.DotaService.GetOrCreate(ctx, dashboardID)
	if err != nil {
		return dotaentity.Nil, fmt.Errorf("failed to get dota settings: %w", err)
	}

	return settings, nil
}
