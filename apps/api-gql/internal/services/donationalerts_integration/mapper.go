package donationalerts_integration

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

// MapEntityToGQLModel converts service entity to GraphQL model
func MapEntityToGQLModel(e entity.DonationAlertsIntegration) *gqlmodel.DonationAlertsIntegration {
	var username, avatar *string

	if e.Data != nil {
		if e.Data.UserName != nil && *e.Data.UserName != "" {
			username = e.Data.UserName
		}
		if e.Data.Avatar != nil && *e.Data.Avatar != "" {
			avatar = e.Data.Avatar
		}
	}

	return &gqlmodel.DonationAlertsIntegration{
		Enabled:  e.Enabled,
		Avatar:   avatar,
		UserName: username,
	}
}
