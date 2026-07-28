package resolvers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlerrors"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	model "github.com/twirapp/twir/libs/gomodels"
)

func (r *mutationResolver) resetCommunityStats(
	ctx context.Context,
	userID string,
	dashboardID string,
	typeArg gqlmodel.CommunityUsersResetType,
) (bool, error) {
	dashboardUUID, err := uuid.Parse(dashboardID)
	if err != nil {
		return false, gqlerrors.HandleError(err)
	}
	if r.deps.DashboardAccess == nil {
		return false, gqlerrors.HandleError(fmt.Errorf("dashboard access service is not configured"))
	}

	isOwner, err := r.deps.DashboardAccess.IsOwner(ctx, userID, dashboardUUID)
	if err != nil {
		return false, gqlerrors.HandleError(err)
	}
	if !isOwner {
		return false, fmt.Errorf("you cannot reset stats for this user")
	}

	var field string

	switch typeArg {
	case gqlmodel.CommunityUsersResetTypeMessages:
		field = "messages"
	case gqlmodel.CommunityUsersResetTypeWatched:
		field = "watched"
	case gqlmodel.CommunityUsersResetTypeUsedChannelsPoints:
		field = "usedChannelPoints"
	case gqlmodel.CommunityUsersResetTypeUsedEmotes:
		field = "emotes"
	}

	if field == "" {
		return false, fmt.Errorf("unknown reset typeArg: %s", typeArg)
	}

	err = r.deps.Gorm.WithContext(ctx).
		Model(&model.UsersStats{}).
		Where(`channel_id = ?`, dashboardID).
		Update(field, 0).Error
	if err != nil {
		return false, gqlerrors.HandleError(err)
	}

	return true, nil
}
