package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func MapChannelFileToGql(m entity.ChannelFile) gqlmodel.ChannelFile {
	return gqlmodel.ChannelFile{
		ID:        m.ID,
		ChannelID: m.ChannelID,
		Mimetype:  m.MimeType,
		Name:      m.FileName,
		Size:      int(m.Size),
	}
}
