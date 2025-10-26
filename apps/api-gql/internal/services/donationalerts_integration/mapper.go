package donationalerts_integration

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

// MapEntityToGQLModel converts service entity to GraphQL model
func MapEntityToGQLModel(e entity.DonationAlertsIntegration) *gqlmodel.DonationAlertsIntegration {
	return &gqlmodel.DonationAlertsIntegration{
		Enabled:  e.Enabled,
		UserName: &e.UserName,
		Avatar:   &e.Avatar,
	}
}
