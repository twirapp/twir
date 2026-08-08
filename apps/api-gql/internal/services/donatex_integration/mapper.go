package donatex_integration

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	donatexintegrationentity "github.com/twirapp/twir/libs/entities/donatex_integration"
)

// MapEntityToGQLModel converts service entity to GraphQL model
func MapEntityToGQLModel(e donatexintegrationentity.Entity) *gqlmodel.DonateXIntegration {
	return &gqlmodel.DonateXIntegration{
		Enabled:  e.Enabled,
		UserName: &e.UserName,
		Avatar:   &e.Avatar,
	}
}
