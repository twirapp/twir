package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/scalars"
	"github.com/twirapp/twir/libs/repositories/channels_storage/model"
)

func StorageEntryModelToGql(m model.ChannelStorage) gqlmodel.StorageEntry {
	return gqlmodel.StorageEntry{
		Key:       m.Key,
		Value:     scalars.JSON(m.Value),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
