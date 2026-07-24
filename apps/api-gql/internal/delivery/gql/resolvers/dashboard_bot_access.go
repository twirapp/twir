package resolvers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
	model "github.com/twirapp/twir/libs/gomodels"
)

const manageBotSettingsPermission = "MANAGE_BOT_SETTINGS"

type dashboardBotUserGetter interface {
	GetAuthenticatedUserModel(context.Context) (*model.Users, error)
}

type dashboardBotAccessChecker interface {
	CanAccess(context.Context, dashboardaccess.Subject, uuid.UUID, string) (bool, error)
}

func ensureBotJoinLeaveAccess(
	ctx context.Context,
	users dashboardBotUserGetter,
	access dashboardBotAccessChecker,
	dashboardID string,
) error {
	dashboardUUID, err := uuid.Parse(dashboardID)
	if err != nil {
		return fmt.Errorf("parse dashboard id: %w", err)
	}

	user, err := users.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return fmt.Errorf("get authenticated user: %w", err)
	}

	hasAccess, err := access.CanAccess(ctx, dashboardaccess.Subject{
		ID:         user.ID,
		IsBotAdmin: user.IsBotAdmin,
	}, dashboardUUID, manageBotSettingsPermission)
	if err != nil {
		return fmt.Errorf("check dashboard access: %w", err)
	}
	if !hasAccess {
		return fmt.Errorf("user has no permission to manage bot settings")
	}

	return nil
}
