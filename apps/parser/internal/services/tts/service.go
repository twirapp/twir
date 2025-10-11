package tts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"github.com/twirapp/twir/apps/parser/locales"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/i18n"
	"github.com/twirapp/twir/libs/repositories/channels_modules_settings_tts"
	"github.com/twirapp/twir/libs/repositories/channels_modules_settings_tts/model"
	"github.com/twirapp/twir/libs/types/types/api/modules"
)

type Service struct {
	repository channels_modules_settings_tts.Repository
	config     *config.Config
}

func New(
	repository channels_modules_settings_tts.Repository,
	config *config.Config,
) *Service {
	return &Service{
		repository: repository,
		config:     config,
	}
}

type Voice struct {
	Name    string
	Country string
}

// GetChannelSettings retrieves TTS settings for a channel
func (s *Service) GetChannelSettings(ctx context.Context, channelID string) (
	*modules.TTSSettings,
	*model.ChannelModulesSettingsTTS,
	error,
) {
	settings, err := s.repository.GetByChannelID(ctx, channelID)
	if err != nil {
		if errors.Is(err, channels_modules_settings_tts.ErrNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	ttsSettings := &modules.TTSSettings{
		Enabled:                            settings.Enabled,
		Rate:                               settings.Rate,
		Volume:                             settings.Volume,
		Pitch:                              settings.Pitch,
		Voice:                              settings.Voice,
		AllowUsersChooseVoiceInMainCommand: settings.AllowUsersChooseVoiceInMainCommand,
		MaxSymbols:                         settings.MaxSymbols,
		DisallowedVoices:                   settings.DisallowedVoices,
		DoNotReadEmoji:                     settings.DoNotReadEmoji,
		DoNotReadTwitchEmotes:              settings.DoNotReadTwitchEmotes,
		DoNotReadLinks:                     settings.DoNotReadLinks,
		ReadChatMessages:                   settings.ReadChatMessages,
		ReadChatMessagesNicknames:          settings.ReadChatMessagesNicknames,
	}

	return ttsSettings, &settings, nil
}

// GetUserSettings retrieves TTS settings for a specific user in a channel
func (s *Service) GetUserSettings(
	ctx context.Context,
	channelID, userID string,
) (*modules.TTSSettings, *model.ChannelModulesSettingsTTS, error) {
	settings, err := s.repository.GetByChannelIDAndUserID(ctx, channelID, userID)
	if err != nil {
		if errors.Is(err, channels_modules_settings_tts.ErrNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	ttsSettings := &modules.TTSSettings{
		Enabled:                            settings.Enabled,
		Rate:                               settings.Rate,
		Volume:                             settings.Volume,
		Pitch:                              settings.Pitch,
		Voice:                              settings.Voice,
		AllowUsersChooseVoiceInMainCommand: settings.AllowUsersChooseVoiceInMainCommand,
		MaxSymbols:                         settings.MaxSymbols,
		DisallowedVoices:                   settings.DisallowedVoices,
		DoNotReadEmoji:                     settings.DoNotReadEmoji,
		DoNotReadTwitchEmotes:              settings.DoNotReadTwitchEmotes,
		DoNotReadLinks:                     settings.DoNotReadLinks,
		ReadChatMessages:                   settings.ReadChatMessages,
		ReadChatMessagesNicknames:          settings.ReadChatMessagesNicknames,
	}

	return ttsSettings, &settings, nil
}

// UpdateChannelSettings updates TTS settings for a channel
func (s *Service) UpdateChannelSettings(
	ctx context.Context,
	channelID string,
	settings *modules.TTSSettings,
) error {
	input := channels_modules_settings_tts.CreateOrUpdateInput{
		ChannelID:                          channelID,
		UserID:                             nil,
		Enabled:                            settings.Enabled,
		Rate:                               settings.Rate,
		Volume:                             settings.Volume,
		Pitch:                              settings.Pitch,
		Voice:                              settings.Voice,
		AllowUsersChooseVoiceInMainCommand: settings.AllowUsersChooseVoiceInMainCommand,
		MaxSymbols:                         settings.MaxSymbols,
		DisallowedVoices:                   settings.DisallowedVoices,
		DoNotReadEmoji:                     settings.DoNotReadEmoji,
		DoNotReadTwitchEmotes:              settings.DoNotReadTwitchEmotes,
		DoNotReadLinks:                     settings.DoNotReadLinks,
		ReadChatMessages:                   settings.ReadChatMessages,
		ReadChatMessagesNicknames:          settings.ReadChatMessagesNicknames,
	}

	_, err := s.repository.UpdateForChannel(ctx, channelID, input)
	return err
}

// UpdateUserSettings updates TTS settings for a specific user in a channel
func (s *Service) UpdateUserSettings(
	ctx context.Context,
	channelID, userID string,
	settings *modules.TTSSettings,
) error {
	input := channels_modules_settings_tts.CreateOrUpdateInput{
		ChannelID:                          channelID,
		UserID:                             nil,
		Enabled:                            settings.Enabled,
		Rate:                               settings.Rate,
		Volume:                             settings.Volume,
		Pitch:                              settings.Pitch,
		Voice:                              settings.Voice,
		AllowUsersChooseVoiceInMainCommand: settings.AllowUsersChooseVoiceInMainCommand,
		MaxSymbols:                         settings.MaxSymbols,
		DisallowedVoices:                   settings.DisallowedVoices,
		DoNotReadEmoji:                     settings.DoNotReadEmoji,
		DoNotReadTwitchEmotes:              settings.DoNotReadTwitchEmotes,
		DoNotReadLinks:                     settings.DoNotReadLinks,
		ReadChatMessages:                   settings.ReadChatMessages,
		ReadChatMessagesNicknames:          settings.ReadChatMessagesNicknames,
	}

	_, err := s.repository.UpdateForUser(ctx, channelID, userID, input)
	return err
}

// CreateUserSettings creates new TTS settings for a user
func (s *Service) CreateUserSettings(
	ctx context.Context,
	channelID, userID string,
	rate, pitch int,
	voice string,
) (*modules.TTSSettings, error) {
	userSettings := &modules.TTSSettings{
		Enabled: lo.ToPtr(true),
		Rate:    rate,
		Volume:  70,
		Pitch:   pitch,
		Voice:   voice,
	}

	input := channels_modules_settings_tts.CreateOrUpdateInput{
		ChannelID:                          channelID,
		UserID:                             &userID,
		Enabled:                            userSettings.Enabled,
		Rate:                               userSettings.Rate,
		Volume:                             userSettings.Volume,
		Pitch:                              userSettings.Pitch,
		Voice:                              userSettings.Voice,
		AllowUsersChooseVoiceInMainCommand: userSettings.AllowUsersChooseVoiceInMainCommand,
		MaxSymbols:                         userSettings.MaxSymbols,
		DisallowedVoices:                   userSettings.DisallowedVoices,
		DoNotReadEmoji:                     userSettings.DoNotReadEmoji,
		DoNotReadTwitchEmotes:              userSettings.DoNotReadTwitchEmotes,
		DoNotReadLinks:                     userSettings.DoNotReadLinks,
		ReadChatMessages:                   userSettings.ReadChatMessages,
		ReadChatMessagesNicknames:          userSettings.ReadChatMessagesNicknames,
	}

	_, err := s.repository.CreateForUser(ctx, input)
	if err != nil {
		return nil, err
	}

	return userSettings, nil
}

// ToggleChannelEnabled toggles the enabled state for a channel's TTS settings
func (s *Service) ToggleChannelEnabled(ctx context.Context, channelID string, enabled bool) error {
	channelSettings, _, err := s.GetChannelSettings(ctx, channelID)
	if err != nil {
		return err
	}

	if channelSettings == nil {
		return errors.New(i18n.GetCtx(ctx, locales.Translations.Services.Tts.Info.NotConfigured))
	}

	channelSettings.Enabled = &enabled
	return s.UpdateChannelSettings(ctx, channelID, channelSettings)
}

// GetAvailableVoices retrieves available TTS voices from the service
func (s *Service) GetAvailableVoices(ctx context.Context) []Voice {
	data := map[string]any{}
	_, err := req.R().
		SetContext(ctx).
		SetSuccessResult(&data).
		Get(fmt.Sprintf("http://%s/info", s.config.TTSServiceUrl))
	if err != nil {
		return nil
	}

	bytes, err := json.Marshal(&data)
	if err != nil {
		return nil
	}

	parsedJson := gjson.ParseBytes(bytes)
	voices := []Voice{}
	parsedJson.Get("rhvoice_wrapper_voices_info").ForEach(
		func(key, value gjson.Result) bool {
			voices = append(
				voices, Voice{
					Name:    key.String(),
					Country: value.Get("country").String(),
				},
			)

			return true
		},
	)

	return voices
}

// ValidateVoice checks if a voice is valid and not disallowed
func (s *Service) ValidateVoice(ctx context.Context, channelID, voiceName string) (Voice, error) {
	voices := s.GetAvailableVoices(ctx)
	if len(voices) == 0 {
		return Voice{}, errors.New(i18n.GetCtx(ctx, locales.Translations.Services.Tts.Info.NoVoices))
	}

	wantedVoice, ok := lo.Find(
		voices, func(item Voice) bool {
			return item.Name == strings.ToLower(voiceName)
		},
	)
	if !ok {
		return Voice{}, fmt.Errorf(i18n.GetCtx(
			ctx,
			locales.Translations.Services.Tts.Errors.NotFound.
				SetVars(locales.KeysServicesTtsErrorsNotFoundVars{UserVoice: voiceName}),
		))
	}

	channelSettings, _, err := s.GetChannelSettings(ctx, channelID)
	if err != nil {
		return Voice{}, err
	}

	if channelSettings != nil {
		_, isDisallowed := lo.Find(
			channelSettings.DisallowedVoices, func(item string) bool {
				return item == wantedVoice.Name
			},
		)

		if isDisallowed {
			return Voice{}, fmt.Errorf(i18n.GetCtx(
				ctx,
				locales.Translations.Services.Tts.Errors.VoiceDisallowed.
					SetVars(locales.KeysServicesTtsErrorsVoiceDisallowedVars{UserVoice: wantedVoice.Name}),
			))
		}
	}

	return wantedVoice, nil
}

// GetFilteredVoices returns voices that are not disallowed for a channel
func (s *Service) GetFilteredVoices(ctx context.Context, channelID string) ([]Voice, error) {
	voices := s.GetAvailableVoices(ctx)
	if len(voices) == 0 {
		return nil, errors.New(i18n.GetCtx(ctx, locales.Translations.Services.Tts.Info.NoVoices))
	}

	channelSettings, _, err := s.GetChannelSettings(ctx, channelID)
	if err != nil {
		return nil, err
	}

	if channelSettings != nil && len(channelSettings.DisallowedVoices) > 0 {
		voices = lo.Filter(
			voices, func(item Voice, _ int) bool {
				return !lo.Contains(channelSettings.DisallowedVoices, item.Name)
			},
		)
	}

	return voices, nil
}

// IsValidURL checks if the input string is a valid URL
func (s *Service) IsValidURL(input string) bool {
	if u, e := url.Parse(input); e == nil {
		if u.Host != "" {
			return s.dnsCheck(u.Host)
		}
		return s.dnsCheck(input)
	}
	return false
}

// dnsCheck performs DNS lookup to validate a host
func (s *Service) dnsCheck(input string) bool {
	input = strings.TrimPrefix(input, "https://")
	input = strings.TrimPrefix(input, "http://")

	ips, _ := net.LookupIP(input)
	return len(ips) > 0
}
