package messagehandler

import (
	"context"
	"errors"
	"strings"

	model "github.com/satont/twir/libs/gomodels"
	"gorm.io/gorm"
)

func (c *MessageHandler) handleTts(ctx context.Context, msg handleMessage) error {
	if strings.HasPrefix(msg.Message.Text, msg.EnrichedData.ChannelCommandPrefix) {
		return nil
	}

	settings, err := c.ttsService.GetChannelTTSSettings(ctx, msg.BroadcasterUserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	if !settings.ReadChatMessages || settings.Enabled == nil || !*settings.Enabled {
		return nil
	}

	ttsCommand := &model.ChannelsCommands{}
	err = c.gorm.WithContext(ctx).
		Where(`"channelId" = ?`, msg.BroadcasterUserId).
		Where(`"module" = ?`, "TTS").
		Where(`"defaultName" = ?`, "tts").
		Find(&ttsCommand).
		Error
	if err != nil {
		return err
	}

	if ttsCommand.ID == "" {
		return nil
	}

	if !ttsCommand.Enabled {
		return nil
	}

	var msgText strings.Builder
	msgText.WriteString(msg.EnrichedData.ChannelCommandPrefix + ttsCommand.Name)

	if settings.ReadChatMessagesNicknames {
		msgText.WriteString(" " + msg.ChatterUserLogin)
	}

	msgText.WriteString(" " + msg.Message.Text)

	text := msgText.String()

	// copy message to avoid changing original message
	newMessage := msg.TwitchChatMessage
	originalCopy := *msg.TwitchChatMessage.Message
	originalCopy.Text = text
	newMessage.Message = &originalCopy

	_, err = c.bus.Parser.ProcessMessageAsCommand.Request(ctx, newMessage)
	if err != nil {
		return err
	}

	return nil
}
