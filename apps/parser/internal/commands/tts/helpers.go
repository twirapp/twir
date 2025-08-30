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
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/repositories/channels_modules_settings_tts"
	"github.com/twirapp/twir/libs/repositories/channels_modules_settings_tts/model"
	"github.com/twirapp/twir/libs/types/types/api/modules"
)

func getSettings(
	ctx context.Context,
	repo channels_modules_settings_tts.Repository,
	channelId, userId string,
) (*modules.TTSSettings, *model.ChannelModulesSettingsTTS, error) {
	var settings model.ChannelModulesSettingsTTS
	var err error

	if userId == channelId || userId == "" {
		// Get channel-level settings
		settings, err = repo.GetByChannelID(ctx, channelId)
	} else {
		// Get user-specific settings
		settings, err = repo.GetByChannelIDAndUserID(ctx, channelId, userId)
	}

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

type Voice struct {
	Name    string
	Country string
}

func getVoices(ctx context.Context, cfg *config.Config) []Voice {
	data := map[string]any{}
	_, err := req.R().
		SetContext(ctx).
		SetSuccessResult(&data).
		Get(fmt.Sprintf("http://%s/info", cfg.TTSServiceUrl))
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

func updateSettings(
	ctx context.Context,
	repo channels_modules_settings_tts.Repository,
	entity *model.ChannelModulesSettingsTTS,
	settings *modules.TTSSettings,
) error {
	input := channels_modules_settings_tts.CreateOrUpdateInput{
		ChannelID:                          entity.ChannelID,
		UserID:                             entity.UserID,
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

	if entity.UserID != nil {
		// Update user-specific settings
		_, err := repo.UpdateForUser(ctx, entity.ChannelID, *entity.UserID, input)
		return err
	} else {
		// Update channel-level settings
		_, err := repo.UpdateForChannel(ctx, entity.ChannelID, input)
		return err
	}
}

func createUserSettings(
	ctx context.Context,
	repo channels_modules_settings_tts.Repository,
	rate, pitch int,
	voice, channelId, userId string,
) (
	*model.ChannelModulesSettingsTTS,
	*modules.TTSSettings,
	error,
) {
	userSettings := &modules.TTSSettings{
		Enabled: lo.ToPtr(true),
		Rate:    rate,
		Volume:  70,
		Pitch:   pitch,
		Voice:   voice,
	}

	input := channels_modules_settings_tts.CreateOrUpdateInput{
		ChannelID:                          channelId,
		UserID:                             &userId,
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

	userModel, err := repo.CreateForUser(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	return &userModel, userSettings, nil
}

func switchEnableState(
	ctx context.Context,
	repo channels_modules_settings_tts.Repository,
	channelId string,
	newState bool,
) error {
	channelSettings, channelModel, err := getSettings(ctx, repo, channelId, "")
	if err != nil {
		return err
	}

	if channelSettings == nil {
		return errors.New("tts not configured")
	}

	channelSettings.Enabled = &newState
	err = updateSettings(ctx, repo, channelModel, channelSettings)
	if err != nil {
		return err
	}

	return nil
}

func isValidUrl(input string) bool {
	if u, e := url.Parse(input); e == nil {
		if u.Host != "" {
			return dnsCheck(u.Host)
		}

		return dnsCheck(input)
	}

	return false
}

func dnsCheck(input string) bool {
	input = strings.TrimPrefix(input, "https://")
	input = strings.TrimPrefix(input, "http://")

	ips, _ := net.LookupIP(input)
	return len(ips) > 0
}
