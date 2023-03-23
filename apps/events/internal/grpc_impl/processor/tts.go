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

func (c *Processor) getTtsSettings(channelId, userId string) (*modules.TTSSettings, *model.ChannelModulesSettings) {
	settings := &model.ChannelModulesSettings{}
	query := c.services.DB.
		Where(`"channelId" = ?`, channelId).
		Where(`"type" = ?`, "tts")

	if userId != "" {
		query = query.Where(`"userId" = ?`, userId)
	} else {
		query = query.Where(`"userId" IS NULL`)
	}

	err := query.First(&settings).Error
	if err != nil {
		return nil, nil
	}

	data := modules.TTSSettings{}
	err = json.Unmarshal(settings.Settings, &data)
	if err != nil {
		return nil, nil
	}

	return &data, settings
}

func (c *Processor) TtsSay(channelId, userId, message string) error {
	msg, err := c.HydrateStringWithData(message, c.data)
	if err != nil {
		return fmt.Errorf("cannot hydrate string %s", err)
	}

	channelSettings, _ := c.getTtsSettings(channelId, "")

	if channelSettings == nil || !*channelSettings.Enabled {
		return InternalError
	}

	userSettings, _ := c.getTtsSettings(channelId, userId)

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

func (c *Processor) TtsSkip(channelId string) error {
	_, err := c.services.WebsocketsGrpc.TextToSpeechSkip(context.Background(), &websockets.TTSSkipMessage{
		ChannelId: channelId,
	})

	if err != nil {
		return fmt.Errorf("cannot send message %s", err)
	}

	return nil
}

func (c *Processor) TtsChangeState(channelId string, enabled bool) error {
	currentSettings, dbModel := c.getTtsSettings(channelId, "")
	if currentSettings == nil {
		return nil
	}

	currentSettings.Enabled = &enabled

	bytes, err := json.Marshal(currentSettings)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return err
	}

	err = c.services.DB.Model(&dbModel).Updates(map[string]interface{}{"settings": bytes}).Error
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return err
	}

	return nil
}
