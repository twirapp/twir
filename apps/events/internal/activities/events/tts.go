package events

import (
	"context"
	"fmt"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/events/internal/shared"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/repositories/channels_modules_settings_tts"
	ttsmodel "github.com/twirapp/twir/libs/repositories/channels_modules_settings_tts/model"
	"github.com/twirapp/twir/libs/repositories/events/model"
	"go.temporal.io/sdk/activity"
)

func (c *Activity) getTtsSettings(
	ctx context.Context,
	channelId,
	userId string,
) (*ttsmodel.ChannelModulesSettingsTTS, error) {
	activity.RecordHeartbeat(ctx, nil)

	if userId != "" {
		data, err := c.ttsRepository.GetByChannelIDAndUserID(ctx, channelId, userId)
		if err != nil {
			return nil, fmt.Errorf("cannot get tts settings by channel id and user id %s", err)
		}

		return &data, nil
	}

	data, err := c.ttsRepository.GetByChannelID(ctx, channelId)
	if err != nil {
		return nil, fmt.Errorf("cannot get tts settings by channel id %s", err)
	}

	return &data, nil
}

func (c *Activity) TtsSay(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	if operation.Input == nil || *operation.Input == "" {
		return fmt.Errorf("input is required for TTS operation")
	}

	msg, hydrateErr := c.hydrator.HydrateStringWithData(data.ChannelID, *operation.Input, data)
	if hydrateErr != nil {
		return fmt.Errorf("cannot hydrate string %s", hydrateErr)
	}

	channelSettings, err := c.getTtsSettings(ctx, data.ChannelID, "")
	if err != nil {
		return fmt.Errorf("cannot get tts settings %s", err)
	}

	if channelSettings == nil || !*channelSettings.Enabled {
		return nil
	}

	userSettings, _ := c.getTtsSettings(ctx, data.ChannelID, data.UserID)

	voice := lo.IfF(
		userSettings != nil, func() string {
			return userSettings.Voice
		},
	).Else(channelSettings.Voice)
	rate := lo.IfF(
		userSettings != nil, func() int {
			return userSettings.Rate
		},
	).Else(channelSettings.Rate)
	pitch := lo.IfF(
		userSettings != nil, func() int {
			return userSettings.Pitch
		},
	).Else(channelSettings.Pitch)

	_, err = c.websocketsGrpc.TextToSpeechSay(
		ctx,
		&websockets.TTSMessage{
			ChannelId: data.ChannelID,
			Text:      msg,
			Voice:     voice,
			Rate:      strconv.Itoa(rate),
			Pitch:     strconv.Itoa(pitch),
			Volume:    strconv.Itoa(channelSettings.Volume),
		},
	)
	if err != nil {
		return fmt.Errorf("cannot say tts message %s", err)
	}

	return nil
}

func (c *Activity) TtsSkip(
	ctx context.Context,
	_ model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	_, err := c.websocketsGrpc.TextToSpeechSkip(
		ctx,
		&websockets.TTSSkipMessage{
			ChannelId: data.ChannelID,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot skip message %s", err)
	}

	return nil
}

func (c *Activity) TtsChangeState(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	currentSettings, err := c.getTtsSettings(ctx, data.ChannelID, "")
	if err != nil {
		return fmt.Errorf("cannot get tts settings %s", err)
	}

	if currentSettings == nil {
		return nil
	}

	var newState bool
	if operation.Type == model.EventOperationTypeTtsEnable {
		newState = true
	} else {
		newState = false
	}

	_, err = c.ttsRepository.UpdateForChannel(
		ctx,
		data.ChannelID,
		channels_modules_settings_tts.CreateOrUpdateInput{
			ChannelID:                          data.ChannelID,
			UserID:                             nil,
			Enabled:                            &newState,
			Rate:                               currentSettings.Rate,
			Volume:                             currentSettings.Volume,
			Pitch:                              currentSettings.Pitch,
			Voice:                              currentSettings.Voice,
			AllowUsersChooseVoiceInMainCommand: currentSettings.AllowUsersChooseVoiceInMainCommand,
			MaxSymbols:                         currentSettings.MaxSymbols,
			DisallowedVoices:                   currentSettings.DisallowedVoices,
			DoNotReadEmoji:                     currentSettings.DoNotReadEmoji,
			DoNotReadTwitchEmotes:              currentSettings.DoNotReadTwitchEmotes,
			DoNotReadLinks:                     currentSettings.DoNotReadLinks,
			ReadChatMessages:                   currentSettings.ReadChatMessages,
			ReadChatMessagesNicknames:          currentSettings.ReadChatMessagesNicknames,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot update tts settings %s", err)
	}

	return nil
}

func (c *Activity) TtsChangeAutoReadState(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	currentSettings, err := c.getTtsSettings(ctx, data.ChannelID, "")
	if err != nil {
		return fmt.Errorf("cannot get tts settings %s", err)
	}
	if currentSettings == nil {
		return nil
	}

	var newState bool
	if operation.Type == model.EventOperationTypeTtsEnableAutoread {
		newState = true
	} else if operation.Type == model.EventOperationTypeTtsDisableAutoread {
		newState = false
	} else if operation.Type == model.EventOperationTypeTtsSwitchAutoread {
		newState = !currentSettings.ReadChatMessages
	}

	_, err = c.ttsRepository.UpdateForChannel(
		ctx,
		data.ChannelID,
		channels_modules_settings_tts.CreateOrUpdateInput{
			ChannelID:                          data.ChannelID,
			UserID:                             nil,
			Enabled:                            currentSettings.Enabled,
			Rate:                               currentSettings.Rate,
			Volume:                             currentSettings.Volume,
			Pitch:                              currentSettings.Pitch,
			Voice:                              currentSettings.Voice,
			AllowUsersChooseVoiceInMainCommand: currentSettings.AllowUsersChooseVoiceInMainCommand,
			MaxSymbols:                         currentSettings.MaxSymbols,
			DisallowedVoices:                   currentSettings.DisallowedVoices,
			DoNotReadEmoji:                     currentSettings.DoNotReadEmoji,
			DoNotReadTwitchEmotes:              currentSettings.DoNotReadTwitchEmotes,
			DoNotReadLinks:                     currentSettings.DoNotReadLinks,
			ReadChatMessages:                   newState,
			ReadChatMessagesNicknames:          currentSettings.ReadChatMessagesNicknames,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot update tts settings %s", err)
	}

	return nil
}
