package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func SongRequestPublicTo(m entity.SongRequestPublic) gqlmodel.SongRequestPublic {
	return gqlmodel.SongRequestPublic{
		Title:           m.Title,
		UserID:          m.UserID,
		CreatedAt:       m.CreatedAt,
		SongLink:        m.SongLink,
		DurationSeconds: m.DurationSeconds,
	}
}
