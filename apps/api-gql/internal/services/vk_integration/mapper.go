package vk_integration

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/libs/entities/vk_integration"
)

// MapEntityToGQLModel converts entity to GraphQL model
func MapEntityToGQLModel(e vk_integration.Entity) *gqlmodel.VKIntegration {
	return &gqlmodel.VKIntegration{
		Enabled:  e.Enabled,
		UserName: &e.UserName,
		Avatar:   &e.Avatar,
	}
}
