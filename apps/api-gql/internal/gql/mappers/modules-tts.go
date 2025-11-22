package tts_mapper

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/modules_tts/model"
)

func DbToEntity(m model.TTS) entity.TtsSettings {
	return entity.TtsSettings{
		ID:                                 m.ID,
		ChannelID:                          m.ChannelID,
		UserID:                             m.UserID,
		CreatedAt:                          m.CreatedAt,
		UpdatedAt:                          m.UpdatedAt,
		Enabled:                            m.Enabled,
		Rate:                               m.Rate,
		Volume:                             m.Volume,
		Pitch:                              m.Pitch,
		Voice:                              m.Voice,
		AllowUsersChooseVoiceInMainCommand: m.AllowUsersChooseVoiceInMainCommand,
		MaxSymbols:                         m.MaxSymbols,
		DisallowedVoices:                   m.DisallowedVoices,
		DoNotReadEmoji:                     m.DoNotReadEmoji,
		DoNotReadTwitchEmotes:              m.DoNotReadTwitchEmotes,
		DoNotReadLinks:                     m.DoNotReadLinks,
		ReadChatMessages:                   m.ReadChatMessages,
		ReadChatMessagesNicknames:          m.ReadChatMessagesNicknames,
	}
}

func EntityToGql(e entity.TtsSettings) *gqlmodel.TtsSettings {
	return &gqlmodel.TtsSettings{
		ID:                                 e.ID,
		ChannelID:                          e.ChannelID,
		Enabled:                            e.Enabled,
		Rate:                               e.Rate,
		Volume:                             e.Volume,
		Pitch:                              e.Pitch,
		Voice:                              e.Voice,
		AllowUsersChooseVoiceInMainCommand: e.AllowUsersChooseVoiceInMainCommand,
		MaxSymbols:                         e.MaxSymbols,
		DisallowedVoices:                   e.DisallowedVoices,
		DoNotReadEmoji:                     e.DoNotReadEmoji,
		DoNotReadTwitchEmotes:              e.DoNotReadTwitchEmotes,
		DoNotReadLinks:                     e.DoNotReadLinks,
		ReadChatMessages:                   e.ReadChatMessages,
		ReadChatMessagesNicknames:          e.ReadChatMessagesNicknames,
		CreatedAt:                          e.CreatedAt,
		UpdatedAt:                          e.UpdatedAt,
	}
}

func EntityToUserGql(e entity.TtsSettings) *gqlmodel.TtsUserSettings {
	if e.UserID == nil {
		return nil
	}

	return &gqlmodel.TtsUserSettings{
		ID:        e.ID,
		ChannelID: e.ChannelID,
		UserID:    *e.UserID,
		Rate:      e.Rate,
		Volume:    e.Volume,
		Pitch:     e.Pitch,
		Voice:     e.Voice,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
