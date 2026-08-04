package resolvers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	authsessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/services/importer"
	apperrors "github.com/twirapp/twir/libs/errors"
	model "github.com/twirapp/twir/libs/gomodels"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
)

type integrationOAuthAttemptCreator interface {
	GetSelectedDashboard(context.Context) (string, error)
	GetAuthenticatedUserModel(context.Context) (*model.Users, error)
	CreateIntegrationOAuthAttempt(context.Context, integrationsmodel.Service, uuid.UUID, uuid.UUID) (string, error)
}

type integrationOAuthAttemptConsumer interface {
	GetSelectedDashboard(context.Context) (string, error)
	GetAuthenticatedUserModel(context.Context) (*model.Users, error)
	ConsumeIntegrationOAuthAttempt(
		context.Context,
		string,
		integrationsmodel.Service,
		uuid.UUID,
		uuid.UUID,
		time.Time,
	) error
}

func createIntegrationOAuthLink(
	ctx context.Context,
	sessions integrationOAuthAttemptCreator,
	service integrationsmodel.Service,
	getLink func(context.Context, string) (string, error),
) (string, error) {
	dashboardID, userID, _, err := integrationOAuthIdentity(ctx, sessions)
	if err != nil {
		return "", err
	}

	state, err := sessions.CreateIntegrationOAuthAttempt(ctx, service, dashboardID, userID)
	if err != nil {
		return "", fmt.Errorf("create %s OAuth attempt: %w", service, err)
	}

	link, err := getLink(ctx, state)
	if err != nil {
		return "", fmt.Errorf("get %s authorization URL: %w", service, err)
	}
	return link, nil
}

func completeIntegrationOAuth(
	ctx context.Context,
	sessions integrationOAuthAttemptConsumer,
	service integrationsmodel.Service,
	code, state string,
	postCode func(context.Context, string, string) error,
) error {
	dashboardID, userID, dashboardIDString, err := integrationOAuthIdentity(ctx, sessions)
	if err != nil {
		return err
	}

	if err := sessions.ConsumeIntegrationOAuthAttempt(
		ctx,
		state,
		service,
		dashboardID,
		userID,
		time.Now(),
	); err != nil {
		if errors.Is(err, authsessions.ErrOAuthAttemptNotFound) ||
			errors.Is(err, authsessions.ErrOAuthAttemptMismatch) ||
			errors.Is(err, authsessions.ErrOAuthAttemptExpired) {
			return apperrors.NewBadRequestError("OAuth authorization attempt is invalid or expired")
		}
		return fmt.Errorf("consume %s OAuth attempt: %w", service, err)
	}

	if err := postCode(ctx, dashboardIDString, code); err != nil {
		return fmt.Errorf("complete %s OAuth authorization: %w", service, err)
	}
	return nil
}

type oauthIdentityReader interface {
	GetSelectedDashboard(context.Context) (string, error)
	GetAuthenticatedUserModel(context.Context) (*model.Users, error)
}

func integrationOAuthIdentity(
	ctx context.Context,
	sessions oauthIdentityReader,
) (uuid.UUID, uuid.UUID, string, error) {
	dashboardIDString, err := sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return uuid.Nil, uuid.Nil, "", fmt.Errorf("get selected dashboard: %w", err)
	}
	dashboardID, err := uuid.Parse(dashboardIDString)
	if err != nil {
		return uuid.Nil, uuid.Nil, "", fmt.Errorf("parse selected dashboard ID: %w", err)
	}

	user, err := sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return uuid.Nil, uuid.Nil, "", fmt.Errorf("get authenticated user: %w", err)
	}
	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return uuid.Nil, uuid.Nil, "", fmt.Errorf("parse authenticated user ID: %w", err)
	}

	return dashboardID, userID, dashboardIDString, nil
}

func mapImportReport(report importer.Report) *gqlmodel.ImportReport {
	failures := make([]gqlmodel.ImportFailure, 0, len(report.Failures))
	for _, failure := range report.Failures {
		failures = append(failures, gqlmodel.ImportFailure{
			Name:   failure.Name,
			Reason: gqlmodel.ImportFailureReason(failure.Reason),
		})
	}

	return &gqlmodel.ImportReport{
		ImportedCount: report.ImportedCount,
		FailedCount:   report.FailedCount,
		Failures:      failures,
	}
}
