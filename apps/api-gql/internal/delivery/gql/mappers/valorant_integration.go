package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func MapValorantIntegrationDataToGql(data *entity.ValorantIntegrationData) *gqlmodel.ValorantIntegration {
	if data == nil {
		return &gqlmodel.ValorantIntegration{
			Enabled: false,
		}
	}

	return &gqlmodel.ValorantIntegration{
		Enabled:  data.Enabled,
		UserName: data.UserName,
		Avatar:   data.Avatar,
	}
}
