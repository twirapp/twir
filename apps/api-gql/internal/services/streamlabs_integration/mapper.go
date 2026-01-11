package streamlabs_integration

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	streamlabsintegration "github.com/twirapp/twir/libs/entities/streamlabs_integration"
)

// MapEntityToGQLModel converts service entity to GraphQL model
func MapEntityToGQLModel(e streamlabsintegration.Entity) *gqlmodel.StreamlabsIntegration {
	return &gqlmodel.StreamlabsIntegration{
		Enabled:  e.Enabled,
		UserName: &e.UserName,
		Avatar:   &e.Avatar,
	}
}
