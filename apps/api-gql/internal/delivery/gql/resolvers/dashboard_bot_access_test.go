package resolvers

import (
	"context"
	"testing"

	"github.com/google/uuid"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
	model "github.com/twirapp/twir/libs/gomodels"
)

type dashboardBotUserGetterStub struct {
	user *model.Users
}

func (s dashboardBotUserGetterStub) GetAuthenticatedUserModel(context.Context) (*model.Users, error) {
	return s.user, nil
}

type dashboardBotAccessCheckerStub struct {
	subject     dashboardaccess.Subject
	dashboardID uuid.UUID
	permission  string
	hasAccess   bool
}

func (s *dashboardBotAccessCheckerStub) CanAccess(
	_ context.Context,
	subject dashboardaccess.Subject,
	dashboardID uuid.UUID,
	permission string,
) (bool, error) {
	s.subject = subject
	s.dashboardID = dashboardID
	s.permission = permission
	return s.hasAccess, nil
}

func TestEnsureBotJoinLeaveAccessChecksExplicitDashboardManagePermission(t *testing.T) {
	// Given
	userID := uuid.New()
	dashboardID := uuid.New()
	access := &dashboardBotAccessCheckerStub{hasAccess: true}
	users := dashboardBotUserGetterStub{user: &model.Users{ID: userID.String()}}

	// When
	err := ensureBotJoinLeaveAccess(context.Background(), users, access, dashboardID.String())

	// Then
	if err != nil {
		t.Fatalf("ensure bot join/leave access: %v", err)
	}
	if access.subject.ID != userID.String() || access.dashboardID != dashboardID {
		t.Fatalf("access check = subject %#v dashboard %s, want explicit dashboard", access.subject, access.dashboardID)
	}
	if access.permission != "MANAGE_BOT_SETTINGS" {
		t.Fatalf("permission = %q, want MANAGE_BOT_SETTINGS", access.permission)
	}
}

func TestEnsureBotJoinLeaveAccessRejectsUnauthorizedExplicitDashboard(t *testing.T) {
	// Given
	access := &dashboardBotAccessCheckerStub{}
	users := dashboardBotUserGetterStub{user: &model.Users{ID: uuid.NewString()}}

	// When
	err := ensureBotJoinLeaveAccess(context.Background(), users, access, uuid.NewString())

	// Then
	if err == nil {
		t.Fatal("ensure bot join/leave access error = nil, want access denial")
	}
}
