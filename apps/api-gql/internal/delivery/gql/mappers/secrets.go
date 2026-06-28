package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/libs/repositories/channels_secret/model"
)

func SecretModelToGql(m model.ChannelSecret) gqlmodel.Secret {
	return gqlmodel.Secret{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description.Ptr(),
	}
}
