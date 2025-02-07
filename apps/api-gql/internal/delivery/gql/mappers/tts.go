package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func TTSUserSettingTo(m entity.TTSUserSettings) gqlmodel.TTSUserSettings {
	return gqlmodel.TTSUserSettings{
		UserID:         m.UserID,
		Rate:           m.Rate,
		Pitch:          m.Pitch,
		Volume:         m.Volume,
		Voice:          m.Voice,
		IsChannelOwner: m.IsChannelOwner,
	}
}
