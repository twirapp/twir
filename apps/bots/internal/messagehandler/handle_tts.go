package messagehandler

import (
	"context"
	"strings"

	"github.com/goccy/go-json"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/types/types/api/modules"
)

func (c *MessageHandler) handleTts(ctx context.Context, msg handleMessage) error {
	if strings.HasPrefix(msg.Message.Text, "!") {
		return nil
	}

	settings := &model.ChannelModulesSettings{}
	query := c.gorm.WithContext(ctx).
		Where(`"channelId" = ?`, msg.BroadcasterUserId).
		Where(`"userId" IS NULL`).
		Where(`"type" = ?`, "tts")

	err := query.Find(&settings).Error
	if err != nil {
		return err
	}

	if settings.ID == "" {
		return nil
	}

	data := modules.TTSSettings{}
	err = json.Unmarshal(settings.Settings, &data)
	if err != nil {
		return err
	}

	if !data.ReadChatMessages || data.Enabled == nil || !*data.Enabled {
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
	msgText.WriteString("!" + ttsCommand.Name)

	if data.ReadChatMessagesNicknames {
		msgText.WriteString(" " + msg.ChatterUserLogin)
	}

	msgText.WriteString(" " + msg.Message.Text)

	text := msgText.String()
	msg.Message.Text = text

	_, err = c.bus.ParserProcessMessageAsCommand.Request(ctx, msg.TwitchChatMessage)
	if err != nil {
		return err
	}

	return nil
}
