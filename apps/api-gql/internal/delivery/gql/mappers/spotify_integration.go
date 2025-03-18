package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func MapSpotifyIntegrationDataToGql(
	data *entity.SpotifyIntegrationData,
) *gqlmodel.SpotifyIntegration {
	if data == nil {
		return nil
	}

	return &gqlmodel.SpotifyIntegration{
		UserName: data.UserName,
		Avatar:   data.Avatar,
	}
}
