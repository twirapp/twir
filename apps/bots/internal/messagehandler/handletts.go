package messagehandler

import (
	"context"
	"strings"

	"github.com/goccy/go-json"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/types/types/api/modules"
	"github.com/twirapp/twir/libs/grpc/parser"
)

func (c *MessageHandler) handleTts(ctx context.Context, msg handleMessage) error {
	if strings.HasPrefix(msg.GetMessage().GetText(), "!") {
		return nil
	}

	settings := &model.ChannelModulesSettings{}
	query := c.gorm.WithContext(ctx).
		Where(`"channelId" = ?`, msg.GetBroadcasterUserId()).
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
		Where(`"channelId" = ?`, msg.GetBroadcasterUserId()).
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
		msgText.WriteString(" " + msg.GetChatterUserLogin())
	}

	msgText.WriteString(" " + msg.GetMessage().GetText())

	text := msgText.String()

	requestStruct := &parser.ProcessCommandRequest{
		Sender: &parser.Sender{
			Id:          msg.GetChatterUserId(),
			Name:        msg.GetChatterUserLogin(),
			DisplayName: msg.GetChatterUserName(),
			Badges:      createUserBadges(msg.GetBadges()),
		},
		Channel: &parser.Channel{
			Id:   msg.GetBroadcasterUserId(),
			Name: msg.GetBroadcasterUserLogin(),
		},
		Message: &parser.Message{
			Id:     msg.GetMessageId(),
			Text:   text,
			Emotes: []*parser.Message_Emote{},
		},
	}

	_, err = c.parserGrpc.ProcessCommand(ctx, requestStruct)
	if err != nil {
		return err
	}

	return nil
}
