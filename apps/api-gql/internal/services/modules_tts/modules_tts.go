package modules_tts

import (
	"context"
	"errors"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	tts_mapper "github.com/twirapp/twir/apps/api-gql/internal/gql/mappers"
	"github.com/twirapp/twir/libs/repositories/modules_tts"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	ModulesTtsRepository modules_tts.Repository
}

func New(opts Opts) *Service {
	return &Service{
		modulesTtsRepository: opts.ModulesTtsRepository,
	}
}

type Service struct {
	modulesTtsRepository modules_tts.Repository
}

type GetSettingsInput struct {
	ChannelID string
}

func (s *Service) GetSettings(ctx context.Context, input GetSettingsInput) (
	entity.TtsSettings,
	error,
) {
	tts, err := s.modulesTtsRepository.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		if errors.Is(err, modules_tts.ErrNotFound) {
			// Create default settings
			return s.CreateSettings(
				ctx, CreateSettingsInput{
					ChannelID:                          input.ChannelID,
					Rate:                               50,
					Volume:                             100,
					Pitch:                              50,
					Voice:                              "",
					AllowUsersChooseVoiceInMainCommand: false,
					MaxSymbols:                         500,
					DisallowedVoices:                   []string{},
					DoNotReadEmoji:                     false,
					DoNotReadTwitchEmotes:              false,
					DoNotReadLinks:                     false,
					ReadChatMessages:                   false,
					ReadChatMessagesNicknames:          false,
				},
			)
		}
		return entity.TtsSettings{}, err
	}

	return tts_mapper.DbToEntity(tts), nil
}

type CreateSettingsInput struct {
	ChannelID                          string
	Enabled                            *bool
	Rate                               int
	Volume                             int
	Pitch                              int
	Voice                              string
	AllowUsersChooseVoiceInMainCommand bool
	MaxSymbols                         int
	DisallowedVoices                   []string
	DoNotReadEmoji                     bool
	DoNotReadTwitchEmotes              bool
	DoNotReadLinks                     bool
	ReadChatMessages                   bool
	ReadChatMessagesNicknames          bool
}

func (s *Service) CreateSettings(
	ctx context.Context,
	input CreateSettingsInput,
) (entity.TtsSettings, error) {
	tts, err := s.modulesTtsRepository.CreateForChannel(
		ctx,
		modules_tts.CreateInput{
			ChannelID:                          input.ChannelID,
			Enabled:                            input.Enabled,
			Rate:                               input.Rate,
			Volume:                             input.Volume,
			Pitch:                              input.Pitch,
			Voice:                              input.Voice,
			AllowUsersChooseVoiceInMainCommand: input.AllowUsersChooseVoiceInMainCommand,
			MaxSymbols:                         input.MaxSymbols,
			DisallowedVoices:                   input.DisallowedVoices,
			DoNotReadEmoji:                     input.DoNotReadEmoji,
			DoNotReadTwitchEmotes:              input.DoNotReadTwitchEmotes,
			DoNotReadLinks:                     input.DoNotReadLinks,
			ReadChatMessages:                   input.ReadChatMessages,
			ReadChatMessagesNicknames:          input.ReadChatMessagesNicknames,
		},
	)
	if err != nil {
		return entity.TtsSettings{}, err
	}

	return tts_mapper.DbToEntity(tts), nil
}

type UpdateSettingsInput struct {
	ChannelID                          string
	Enabled                            *bool
	Rate                               int
	Volume                             int
	Pitch                              int
	Voice                              string
	AllowUsersChooseVoiceInMainCommand bool
	MaxSymbols                         int
	DisallowedVoices                   []string
	DoNotReadEmoji                     bool
	DoNotReadTwitchEmotes              bool
	DoNotReadLinks                     bool
	ReadChatMessages                   bool
	ReadChatMessagesNicknames          bool
}

func (s *Service) UpdateSettings(
	ctx context.Context,
	input UpdateSettingsInput,
) (entity.TtsSettings, error) {
	tts, err := s.modulesTtsRepository.UpdateForChannel(
		ctx,
		input.ChannelID,
		modules_tts.UpdateInput{
			Enabled:                            input.Enabled,
			Rate:                               input.Rate,
			Volume:                             input.Volume,
			Pitch:                              input.Pitch,
			Voice:                              input.Voice,
			AllowUsersChooseVoiceInMainCommand: input.AllowUsersChooseVoiceInMainCommand,
			MaxSymbols:                         input.MaxSymbols,
			DisallowedVoices:                   input.DisallowedVoices,
			DoNotReadEmoji:                     input.DoNotReadEmoji,
			DoNotReadTwitchEmotes:              input.DoNotReadTwitchEmotes,
			DoNotReadLinks:                     input.DoNotReadLinks,
			ReadChatMessages:                   input.ReadChatMessages,
			ReadChatMessagesNicknames:          input.ReadChatMessagesNicknames,
		},
	)
	if err != nil {
		if errors.Is(err, modules_tts.ErrNotFound) {
			// Create if not exists
			return s.CreateSettings(
				ctx, CreateSettingsInput{
					ChannelID:                          input.ChannelID,
					Enabled:                            input.Enabled,
					Rate:                               input.Rate,
					Volume:                             input.Volume,
					Pitch:                              input.Pitch,
					Voice:                              input.Voice,
					AllowUsersChooseVoiceInMainCommand: input.AllowUsersChooseVoiceInMainCommand,
					MaxSymbols:                         input.MaxSymbols,
					DisallowedVoices:                   input.DisallowedVoices,
					DoNotReadEmoji:                     input.DoNotReadEmoji,
					DoNotReadTwitchEmotes:              input.DoNotReadTwitchEmotes,
					DoNotReadLinks:                     input.DoNotReadLinks,
					ReadChatMessages:                   input.ReadChatMessages,
					ReadChatMessagesNicknames:          input.ReadChatMessagesNicknames,
				},
			)
		}
		return entity.TtsSettings{}, err
	}

	return tts_mapper.DbToEntity(tts), nil
}

type GetUserSettingsInput struct {
	ChannelID string
}

func (s *Service) GetUsersSettings(
	ctx context.Context,
	input GetUserSettingsInput,
) ([]entity.TtsSettings, error) {
	users, err := s.modulesTtsRepository.GetAllUsersByChannelID(ctx, input.ChannelID)
	if err != nil {
		return nil, err
	}

	result := make([]entity.TtsSettings, len(users))
	for i, user := range users {
		result[i] = tts_mapper.DbToEntity(user)
	}

	return result, nil
}

type DeleteUsersInput struct {
	ChannelID string
	UserIDs   []string
}

func (s *Service) DeleteUsers(ctx context.Context, input DeleteUsersInput) error {
	return s.modulesTtsRepository.DeleteUsersForChannel(ctx, input.ChannelID, input.UserIDs)
}
