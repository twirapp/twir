package tts

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/imroc/req/v3"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/wsrouter"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/api"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/repositories/overlays_tts"
	ttsmodel "github.com/twirapp/twir/libs/repositories/overlays_tts/model"
	"github.com/twirapp/twir/libs/repositories/users"
	"github.com/twirapp/twir/libs/types/types/api/modules"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Repository      overlays_tts.Repository
	WsRouter        wsrouter.WsRouter
	Config          config.Config
	TwirBus         *buscore.Bus
	Logger          *slog.Logger
	UsersRepository users.Repository
	Cacher          *generic_cacher.GenericCacher[modules.TTSSettings]
}

func New(opts Opts) *Service {
	s := &Service{
		repository:      opts.Repository,
		wsRouter:        opts.WsRouter,
		config:          opts.Config,
		twirBus:         opts.TwirBus,
		usersRepository: opts.UsersRepository,
		cacher:          opts.Cacher,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				s.twirBus.Api.TriggerTtsSay.SubscribeGroup(
					"api",
					func(ctx context.Context, data api.TriggerTtsSay) (struct{}, error) {
						return struct{}{}, s.wsRouter.Publish(createSaySubscriptionKey(data.ChannelId), data)
					},
				)

				opts.Logger.Info("Subscribed to TriggerTtsSay events")

				s.twirBus.Api.TriggerTtsSkip.SubscribeGroup(
					"api",
					func(ctx context.Context, data api.TriggerTtsSkip) (struct{}, error) {
						return struct{}{}, s.wsRouter.Publish(createSkipSubscriptionKey(data.ChannelId), data)
					},
				)

				opts.Logger.Info("Subscribed to TriggerTtsSkip events")

				return nil
			},
			OnStop: func(ctx context.Context) error {
				s.twirBus.Api.TriggerTtsSay.Unsubscribe()
				s.twirBus.Api.TriggerTtsSkip.Unsubscribe()

				opts.Logger.Info("Unsubscribed from TriggerTtsSay and TriggerTtsSkip events")

				return nil
			},
		},
	)

	return s
}

type Service struct {
	repository      overlays_tts.Repository
	wsRouter        wsrouter.WsRouter
	config          config.Config
	twirBus         *buscore.Bus
	usersRepository users.Repository
	cacher          *generic_cacher.GenericCacher[modules.TTSSettings]
}

// GetOrCreate gets the TTS overlay for the given channel ID or creates a new one with default settings if it doesn't exist
func (s *Service) GetOrCreate(ctx context.Context, channelID string) (
	entity.TTSOverlay,
	error,
) {
	overlay, err := s.repository.GetByChannelID(ctx, channelID)
	if err != nil {
		if errors.Is(err, overlays_tts.ErrNotFound) {
			// Create a new overlay with default settings
			created, err := s.repository.Create(ctx, createDefaultOverlayInput(channelID))
			if err != nil {
				return entity.TTSOverlay{}, err
			}
			return mapModelToEntity(created), nil
		}
		return entity.TTSOverlay{}, err
	}

	return mapModelToEntity(overlay), nil
}

type UpdateInput struct {
	ChannelID string
	Settings  entity.TTSOverlaySettings
}

// Update updates the TTS overlay for the given channel ID
func (s *Service) Update(
	ctx context.Context,
	input UpdateInput,
) (entity.TTSOverlay, error) {
	updated, err := s.repository.Update(
		ctx, input.ChannelID, overlays_tts.UpdateInput{
			Settings: mapSettingsEntityToModel(input.Settings),
		},
	)
	if err != nil {
		return entity.TTSOverlay{}, err
	}

	if err := s.wsRouter.Publish(
		createSettingsSubscriptionKey(input.ChannelID),
		mapModelToEntity(updated),
	); err != nil {
		return entity.TTSOverlay{}, err
	}

	if err := s.cacher.Invalidate(ctx, input.ChannelID); err != nil {
		return entity.TTSOverlay{}, fmt.Errorf("invalidate cache: %w", err)
	}

	return mapModelToEntity(updated), nil
}

func (s *Service) GetTTSUsersSettings(
	ctx context.Context,
	channelID string,
) ([]entity.TTSUserSettings, error) {
	userSettings, err := s.repository.GetAllUserSettings(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("get all user settings: %w", err)
	}

	mappedEntities := make([]entity.TTSUserSettings, 0, len(userSettings)+1)

	// Add channel owner settings if they exist
	channelSettings, err := s.repository.GetByChannelID(ctx, channelID)
	if err != nil && !errors.Is(err, overlays_tts.ErrNotFound) {
		return nil, fmt.Errorf("get channel settings: %w", err)
	}

	if err == nil && channelSettings.Settings != nil {
		mappedEntities = append(
			mappedEntities,
			entity.TTSUserSettings{
				UserID:         channelID,
				Rate:           int(channelSettings.Settings.Rate),
				Pitch:          int(channelSettings.Settings.Pitch),
				Voice:          channelSettings.Settings.Voice,
				IsChannelOwner: true,
			},
		)
	}

	// Add user-specific settings
	for _, setting := range userSettings {
		mappedEntities = append(
			mappedEntities,
			entity.TTSUserSettings{
				UserID:         setting.UserID,
				Rate:           int(setting.Rate),
				Pitch:          int(setting.Pitch),
				Voice:          setting.Voice,
				IsChannelOwner: false,
			},
		)
	}

	return mappedEntities, nil
}

func (s *Service) DeleteUsersSettings(
	ctx context.Context,
	channelID string,
	userIds []string,
) error {
	return s.repository.DeleteMultipleUserSettings(ctx, channelID, userIds)
}

// Mappers between repository model and entity
func mapModelToEntity(m ttsmodel.TTSOverlay) entity.TTSOverlay {
	return entity.TTSOverlay{
		ID:        m.ID,
		ChannelID: m.ChannelID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Settings: entity.TTSOverlaySettings{
			Enabled:                            m.Settings.Enabled,
			Voice:                              m.Settings.Voice,
			DisallowedVoices:                   m.Settings.DisallowedVoices,
			Pitch:                              m.Settings.Pitch,
			Rate:                               m.Settings.Rate,
			Volume:                             m.Settings.Volume,
			DoNotReadTwitchEmotes:              m.Settings.DoNotReadTwitchEmotes,
			DoNotReadEmoji:                     m.Settings.DoNotReadEmoji,
			DoNotReadLinks:                     m.Settings.DoNotReadLinks,
			AllowUsersChooseVoiceInMainCommand: m.Settings.AllowUsersChooseVoiceInMainCommand,
			MaxSymbols:                         m.Settings.MaxSymbols,
			ReadChatMessages:                   m.Settings.ReadChatMessages,
			ReadChatMessagesNicknames:          m.Settings.ReadChatMessagesNicknames,
		},
	}
}

func mapSettingsEntityToModel(e entity.TTSOverlaySettings) ttsmodel.TTSOverlaySettings {
	return ttsmodel.TTSOverlaySettings{
		Enabled:                            e.Enabled,
		Voice:                              e.Voice,
		DisallowedVoices:                   pq.StringArray(e.DisallowedVoices),
		Pitch:                              e.Pitch,
		Rate:                               e.Rate,
		Volume:                             e.Volume,
		DoNotReadTwitchEmotes:              e.DoNotReadTwitchEmotes,
		DoNotReadEmoji:                     e.DoNotReadEmoji,
		DoNotReadLinks:                     e.DoNotReadLinks,
		AllowUsersChooseVoiceInMainCommand: e.AllowUsersChooseVoiceInMainCommand,
		MaxSymbols:                         e.MaxSymbols,
		ReadChatMessages:                   e.ReadChatMessages,
		ReadChatMessagesNicknames:          e.ReadChatMessagesNicknames,
	}
}

func createDefaultOverlayInput(channelID string) overlays_tts.CreateInput {
	return overlays_tts.CreateInput{
		ChannelID: channelID,
		Settings: ttsmodel.TTSOverlaySettings{
			Enabled:                            false,
			Voice:                              "alan",
			DisallowedVoices:                   pq.StringArray{},
			Pitch:                              50,
			Rate:                               50,
			Volume:                             30,
			DoNotReadTwitchEmotes:              true,
			DoNotReadEmoji:                     true,
			DoNotReadLinks:                     true,
			AllowUsersChooseVoiceInMainCommand: false,
			MaxSymbols:                         0,
			ReadChatMessages:                   false,
			ReadChatMessagesNicknames:          false,
		},
	}
}

func createSettingsSubscriptionKey(channelID string) string {
	return "tts:settings:" + channelID
}

type VoiceInfo struct {
	Country string `json:"country"`
	Gender  string `json:"gender"`
	Lang    string `json:"lang"`
	Name    string `json:"name"`
	No      int    `json:"no"`
}

type TTSInfo struct {
	VoicesInfo map[string]VoiceInfo
}

func (s *Service) GetInfo(ctx context.Context) (*TTSInfo, error) {
	result := map[string]any{}
	resp, err := req.
		R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(fmt.Sprintf("http://%s/info", s.config.TTSServiceUrl))
	if err != nil {
		return nil, fmt.Errorf("tts service is not available: %w", err)
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("tts service returned error status: %s", resp.Status)
	}

	respVoicesInfo, ok := result["rhvoice_wrapper_voices_info"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("voices info not found in response")
	}

	voicesInfo := make(map[string]VoiceInfo, len(respVoicesInfo))
	for key, value := range respVoicesInfo {
		info, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		voiceInfo := VoiceInfo{}
		if country, ok := info["country"].(string); ok {
			voiceInfo.Country = country
		}
		if gender, ok := info["gender"].(string); ok {
			voiceInfo.Gender = gender
		}
		if lang, ok := info["lang"].(string); ok {
			voiceInfo.Lang = lang
		}
		if name, ok := info["name"].(string); ok {
			voiceInfo.Name = name
		}
		if no, ok := info["no"].(float64); ok {
			voiceInfo.No = int(no)
		}

		voicesInfo[key] = voiceInfo
	}

	return &TTSInfo{
		VoicesInfo: voicesInfo,
	}, nil
}
