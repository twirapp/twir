package handlers

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"github.com/satont/twir/libs/types/types/api/modules"
)

func (c *Handlers) handleTts(msg *Message, userBadges []string) {
	if strings.HasPrefix(msg.Message, "!") {
		return
	}

	settings := &model.ChannelModulesSettings{}
	query := c.db.
		Where(`"channelId" = ?`, msg.Channel.ID).
		Where(`"userId" IS NULL`).
		Where(`"type" = ?`, "tts")

	err := query.Find(&settings).Error
	if err != nil {
		c.logger.Sugar().Error(err)
		return
	}

	if settings.ID == "" {
		return
	}

	data := modules.TTSSettings{}
	err = json.Unmarshal(settings.Settings, &data)
	if err != nil {
		c.logger.Sugar().Error(err)
		return
	}

	if !data.ReadChatMessages || data.Enabled == nil || !*data.Enabled {
		return
	}

	ttsCommand := &model.ChannelsCommands{}
	err = c.db.
		Where(`"channelId" = ?`, msg.Channel.ID).
		Where(`"module" = ?`, "TTS").
		Where(`"defaultName" = ?`, "tts").
		Find(&ttsCommand).
		Error
	if err != nil {
		c.logger.Sugar().Error(err)
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
			Text: msgText.String(),
			Emotes: lo.Map(
				msg.Emotes, func(item *twitch.Emote, _ int) *parser.Message_Emote {
					return &parser.Message_Emote{
						Name:  item.Name,
						Id:    item.ID,
						Count: int64(item.Count),
						Positions: lo.Map(
							item.Positions, func(item twitch.EmotePosition, _ int) *parser.Message_EmotePosition {
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

	_, err = c.parserGrpc.ProcessCommand(context.Background(), requestStruct)
	if err != nil {
		c.logger.Sugar().Error(err)
	}
}
