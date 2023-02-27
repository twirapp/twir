package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"
	"github.com/satont/tsuwari/libs/types/types/api/modules"
	"strconv"
)

func (c *Processor) getTtsSettings(channelId, userId string) *modules.TTSSettings {
	settings := &model.ChannelModulesSettings{}
	query := c.services.DB.
		Where(`"channelId" = ?`, channelId).
		Where(`"type" = ?`, "tts")

	if userId != "" {
		query = query.Where(`"userId" = ?`, userId)
	}

	err := query.First(&settings).Error
	if err != nil {
		return nil
	}

	data := modules.TTSSettings{}
	err = json.Unmarshal(settings.Settings, &data)
	if err != nil {
		return nil
	}

	return &data
}

func (c *Processor) TtsSay(channelId, userId, message string) error {
	msg, err := hydrateStringWithData(message, c.data)
	if err != nil {
		return fmt.Errorf("cannot hydrate string %s", err)
	}

	channelSettings := c.getTtsSettings(channelId, "")

	if channelSettings == nil || !*channelSettings.Enabled {
		return InternalError
	}

	userSettings := c.getTtsSettings(channelId, userId)

	voice := lo.IfF(userSettings != nil, func() string {
		return userSettings.Voice
	}).Else(channelSettings.Voice)
	rate := lo.IfF(userSettings != nil, func() int {
		return userSettings.Rate
	}).Else(channelSettings.Rate)
	pitch := lo.IfF(userSettings != nil, func() int {
		return userSettings.Pitch
	}).Else(channelSettings.Pitch)

	_, err = c.services.WebsocketsGrpc.TextToSpeechSay(context.Background(), &websockets.TTSMessage{
		ChannelId: channelId,
		Text:      msg,
		Voice:     voice,
		Rate:      strconv.Itoa(rate),
		Pitch:     strconv.Itoa(pitch),
		Volume:    strconv.Itoa(channelSettings.Volume),
	})

	if err != nil {
		return fmt.Errorf("cannot send message %s", err)
	}

	return nil
}
