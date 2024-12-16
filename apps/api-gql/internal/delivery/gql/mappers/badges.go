package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func BadgeEntityToGql(m entity.Badge) gqlmodel.Badge {
	return gqlmodel.Badge{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt.String(),
		FileURL:   m.FileURL,
		Enabled:   m.Enabled,
		FfzSlot:   m.FFZSlot,
	}
}
