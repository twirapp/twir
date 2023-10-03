package chat_client

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"github.com/satont/twir/libs/types/types/api/modules"
)

func (c *ChatClient) handleTts(msg *Message, userBadges []string) {
	if strings.HasPrefix(msg.Message, "!") {
		return
	}

	settings := &model.ChannelModulesSettings{}
	query := c.services.DB.
		Where(`"channelId" = ?`, msg.Channel.ID).
		Where(`"userId" IS NULL`).
		Where(`"type" = ?`, "tts")

	err := query.Find(&settings).Error
	if err != nil {
		c.services.Logger.Error(
			"cannot find tts settings",
			slog.Any("err", err),
			slog.String("channelId", msg.Channel.ID),
		)
		return
	}

	if settings.ID == "" {
		return
	}

	data := modules.TTSSettings{}
	err = json.Unmarshal(settings.Settings, &data)
	if err != nil {
		c.services.Logger.Error(
			"cannot unmarshall tts settings",
			slog.Any("err", err),
			slog.String("channelId", msg.Channel.ID),
		)
		return
	}

	if !data.ReadChatMessages || data.Enabled == nil || !*data.Enabled {
		return
	}

	ttsCommand := &model.ChannelsCommands{}
	err = c.services.DB.
		Where(`"channelId" = ?`, msg.Channel.ID).
		Where(`"module" = ?`, "TTS").
		Where(`"defaultName" = ?`, "tts").
		Find(&ttsCommand).
		Error
	if err != nil {
		c.services.Logger.Error(
			"cannot find tts command",
			slog.Any("err", err),
			slog.String("channelId", msg.Channel.ID),
		)
		return
	}

	if ttsCommand.ID == "" {
		return
	}

	if !ttsCommand.Enabled {
		return
	}

	var msgText strings.Builder
	msgText.WriteString("!" + ttsCommand.Name)

	if data.ReadChatMessagesNicknames {
		msgText.WriteString(" " + msg.User.Name)
	}

	msgText.WriteString(" " + msg.Message)

	text := msgText.String()

	requestStruct := &parser.ProcessCommandRequest{
		Sender: &parser.Sender{
			Id:          msg.User.ID,
			Name:        msg.User.Name,
			DisplayName: msg.User.DisplayName,
			Badges:      userBadges,
		},
		Channel: &parser.Channel{
			Id:   msg.Channel.ID,
			Name: msg.Channel.Name,
		},
		Message: &parser.Message{
			Id:   msg.ID,
			Text: text,
			Emotes: lo.Map(
				msg.Emotes,
				func(item MessageEmote, _ int) *parser.Message_Emote {
					return &parser.Message_Emote{
						Name:  item.Name,
						Id:    item.ID,
						Count: int64(item.Count),
						Positions: lo.Map(
							item.Positions, func(item EmotePosition, _ int) *parser.Message_EmotePosition {
								return &parser.Message_EmotePosition{
									Start: int64(item.Start),
									End:   int64(item.End),
								}
							},
						),
					}
				},
			),
		},
	}

	_, err = c.services.ParserGrpc.ProcessCommand(context.Background(), requestStruct)
	if err != nil {
		c.services.Logger.Error(
			"cannot process tts", slog.Any("err", err),
			slog.String("channelId", msg.Channel.ID),
			slog.String("message", text),
		)
	}
}
