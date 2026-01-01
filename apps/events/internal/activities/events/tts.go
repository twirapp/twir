package events

import (
	"context"
	"fmt"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/events/internal/shared"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/repositories/events/model"
	ttsrepository "github.com/twirapp/twir/libs/repositories/overlays_tts"
	ttsmodel "github.com/twirapp/twir/libs/repositories/overlays_tts/model"
	"go.temporal.io/sdk/activity"
)

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

	channelSettings, err := c.ttsCache.Get(ctx, data.ChannelID)
	if err != nil {
		return fmt.Errorf("cannot get tts settings %s", err)
	}

	if channelSettings.Enabled == nil || !*channelSettings.Enabled {
		return nil
	}

	userSettings, _ := c.ttsRepository.GetUserSettings(ctx, data.ChannelID, data.UserID)

	voice := lo.IfF(
		userSettings != ttsmodel.NilUserSettings, func() string {
			return userSettings.Voice
		},
	).Else(channelSettings.Voice)
	rate := lo.IfF(
		userSettings != ttsmodel.NilUserSettings, func() int {
			return int(userSettings.Rate)
		},
	).Else(channelSettings.Rate)
	pitch := lo.IfF(
		userSettings != ttsmodel.NilUserSettings, func() int {
			return int(userSettings.Pitch)
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

	currentSettings, err := c.ttsCache.Get(ctx, data.ChannelID)
	if err != nil {
		return fmt.Errorf("cannot get tts settings %s", err)
	}

	var newState bool
	if operation.Type == model.EventOperationTypeTtsEnable {
		newState = true
	} else {
		newState = false
	}

	_, err = c.ttsRepository.Update(
		ctx,
		data.ChannelID,
		ttsrepository.UpdateInput{
			Settings: ttsmodel.TTSOverlaySettings{
				Enabled:                            newState,
				Rate:                               int32(currentSettings.Rate),
				Volume:                             int32(currentSettings.Volume),
				Pitch:                              int32(currentSettings.Pitch),
				Voice:                              currentSettings.Voice,
				AllowUsersChooseVoiceInMainCommand: currentSettings.AllowUsersChooseVoiceInMainCommand,
				MaxSymbols:                         int32(currentSettings.MaxSymbols),
				DisallowedVoices:                   currentSettings.DisallowedVoices,
				DoNotReadEmoji:                     currentSettings.DoNotReadEmoji,
				DoNotReadTwitchEmotes:              currentSettings.DoNotReadTwitchEmotes,
				DoNotReadLinks:                     currentSettings.DoNotReadLinks,
				ReadChatMessages:                   currentSettings.ReadChatMessages,
				ReadChatMessagesNicknames:          currentSettings.ReadChatMessagesNicknames,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("cannot update tts settings %s", err)
	}

	if err := c.ttsCache.Invalidate(ctx, data.ChannelID); err != nil {
		return fmt.Errorf("cannot invalidate tts cache %s", err)
	}

	return nil
}

func (c *Activity) TtsChangeAutoReadState(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	currentSettings, err := c.ttsCache.Get(ctx, data.ChannelID)
	if err != nil {
		return fmt.Errorf("cannot get tts settings %s", err)
	}

	var newState bool
	if operation.Type == model.EventOperationTypeTtsEnableAutoread {
		newState = true
	} else if operation.Type == model.EventOperationTypeTtsDisableAutoread {
		newState = false
	} else if operation.Type == model.EventOperationTypeTtsSwitchAutoread {
		newState = !currentSettings.ReadChatMessages
	}

	_, err = c.ttsRepository.Update(
		ctx,
		data.ChannelID,
		ttsrepository.UpdateInput{
			Settings: ttsmodel.TTSOverlaySettings{
				Enabled:                            currentSettings.Enabled != nil && *currentSettings.Enabled,
				Rate:                               int32(currentSettings.Rate),
				Volume:                             int32(currentSettings.Volume),
				Pitch:                              int32(currentSettings.Pitch),
				Voice:                              currentSettings.Voice,
				AllowUsersChooseVoiceInMainCommand: currentSettings.AllowUsersChooseVoiceInMainCommand,
				MaxSymbols:                         int32(currentSettings.MaxSymbols),
				DisallowedVoices:                   currentSettings.DisallowedVoices,
				DoNotReadEmoji:                     currentSettings.DoNotReadEmoji,
				DoNotReadTwitchEmotes:              currentSettings.DoNotReadTwitchEmotes,
				DoNotReadLinks:                     currentSettings.DoNotReadLinks,
				ReadChatMessages:                   newState,
				ReadChatMessagesNicknames:          currentSettings.ReadChatMessagesNicknames,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("cannot update tts settings %s", err)
	}

	if err := c.ttsCache.Invalidate(ctx, data.ChannelID); err != nil {
		return fmt.Errorf("cannot invalidate tts cache %s", err)
	}

	return nil
}
