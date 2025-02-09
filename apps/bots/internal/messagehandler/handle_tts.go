package messagehandler

import (
	"context"
	"strings"

	model "github.com/satont/twir/libs/gomodels"
)

func (c *MessageHandler) handleTts(ctx context.Context, msg handleMessage) error {
	commandsPrefix, err := c.commandsService.GetCommandsPrefix(ctx, msg.BroadcasterUserId)
	if err != nil {
		return nil
	}

	if strings.HasPrefix(msg.Message.Text, commandsPrefix) {
		return nil
	}

	settings, err := c.ttsService.GetChannelTTSSettings(ctx, msg.BroadcasterUserId)
	if err != nil {
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
	msgText.WriteString(commandsPrefix + ttsCommand.Name)

	if settings.ReadChatMessagesNicknames {
		msgText.WriteString(" " + msg.ChatterUserLogin)
	}

	msgText.WriteString(" " + msg.Message.Text)

	text := msgText.String()
	msg.Message.Text = text

	_, err = c.bus.Parser.ProcessMessageAsCommand.Request(ctx, msg.TwitchChatMessage)
	if err != nil {
		return err
	}

	return nil
}
