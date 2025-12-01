package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func MapLastfmIntegrationDataToGql(data *entity.LastfmIntegrationData) *gqlmodel.LastfmIntegration {
	if data == nil {
		return &gqlmodel.LastfmIntegration{
			Enabled: false,
		}
	}

	return &gqlmodel.LastfmIntegration{
		Enabled:  data.Enabled,
		UserName: data.UserName,
		Avatar:   data.Avatar,
	}
}
