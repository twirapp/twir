package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func MapTTSOverlayEntityToGQL(e entity.TTSOverlay) gqlmodel.TTSOverlay {
	return gqlmodel.TTSOverlay{
		ID:                                 e.ID,
		Enabled:                            e.Settings.Enabled,
		Voice:                              e.Settings.Voice,
		DisallowedVoices:                   e.Settings.DisallowedVoices,
		Pitch:                              int(e.Settings.Pitch),
		Rate:                               int(e.Settings.Rate),
		Volume:                             int(e.Settings.Volume),
		DoNotReadTwitchEmotes:              e.Settings.DoNotReadTwitchEmotes,
		DoNotReadEmoji:                     e.Settings.DoNotReadEmoji,
		DoNotReadLinks:                     e.Settings.DoNotReadLinks,
		AllowUsersChooseVoiceInMainCommand: e.Settings.AllowUsersChooseVoiceInMainCommand,
		MaxSymbols:                         int(e.Settings.MaxSymbols),
		ReadChatMessages:                   e.Settings.ReadChatMessages,
		ReadChatMessagesNicknames:          e.Settings.ReadChatMessagesNicknames,
		CreatedAt:                          e.CreatedAt,
		UpdatedAt:                          e.UpdatedAt,
		ChannelID:                          e.ChannelID,
	}
}

func TTSUserSettingTo(e entity.TTSUserSettings) gqlmodel.TTSUserSettings {
	return gqlmodel.TTSUserSettings{
		UserID:         e.UserID,
		Rate:           e.Rate,
		Pitch:          e.Pitch,
		Voice:          e.Voice,
		IsChannelOwner: e.IsChannelOwner,
	}
}
