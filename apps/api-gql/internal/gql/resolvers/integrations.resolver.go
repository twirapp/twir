package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"

	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
)

// IntegrationsGetData is the resolver for the integrationsGetData field.
func (r *queryResolver) IntegrationsGetData(
	ctx context.Context,
	service gqlmodel.IntegrationService,
) (gqlmodel.IntegrationData, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	return r.integrationsDataFetcher.GetIntegrationData(ctx, dashboardId, service)
}

// IntegrationsPostCode is the resolver for the integrationsPostCode field.
func (r *queryResolver) IntegrationsPostCode(
	ctx context.Context,
	service gqlmodel.IntegrationService,
	code string,
) (gqlmodel.IntegrationData, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	if err := r.integrationsPostCodeHandler.PostCode(ctx, service, dashboardId, code); err != nil {
		return nil, err
	}

	return r.integrationsDataFetcher.GetIntegrationData(ctx, dashboardId, service)
}
