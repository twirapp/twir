package events

import (
	"context"
	"fmt"
	"strconv"

	"github.com/goccy/go-json"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/events/internal/shared"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/types/types/api/modules"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"go.temporal.io/sdk/activity"
)

func (c *Activity) getTtsSettings(
	ctx context.Context,
	channelId,
	userId string,
) (
	*modules.TTSSettings,
	*model.ChannelModulesSettings,
) {
	activity.RecordHeartbeat(ctx, nil)

	settings := &model.ChannelModulesSettings{}
	query := c.db.
		WithContext(ctx).
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

func (c *Activity) TtsSay(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	if operation.Input.String == "" {
		return nil
	}

	msg, hydrateErr := c.hydrator.HydrateStringWithData(data.ChannelID, operation.Input.String, data)
	if hydrateErr != nil {
		return fmt.Errorf("cannot hydrate string %s", hydrateErr)
	}

	channelSettings, _ := c.getTtsSettings(ctx, data.ChannelID, "")

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

	_, err := c.websocketsGrpc.TextToSpeechSay(
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
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	_, err := c.websocketsGrpc.TextToSpeechSkip(
		ctx, &websockets.TTSSkipMessage{
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

	currentSettings, dbModel := c.getTtsSettings(ctx, data.ChannelID, "")
	if currentSettings == nil {
		return nil
	}

	if operation.Type == model.OperationTTSEnable {
		currentSettings.Enabled = lo.ToPtr(true)
	} else {
		currentSettings.Enabled = lo.ToPtr(false)
	}

	bytes, err := json.Marshal(currentSettings)
	if err != nil {
		return err
	}

	err = c.db.Model(&dbModel).Updates(map[string]interface{}{"settings": bytes}).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *Activity) TtsChangeAutoReadState(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	currentSettings, dbModel := c.getTtsSettings(ctx, data.ChannelID, "")
	if currentSettings == nil {
		return nil
	}

	var newState bool
	if operation.Type == model.OperationTTSEnableAutoRead {
		newState = true
	} else if operation.Type == model.OperationTTSDisableAutoRead {
		newState = false
	} else if operation.Type == model.OperationTTSSwitchAutoRead {
		newState = !currentSettings.ReadChatMessages
	}

	currentSettings.ReadChatMessages = newState

	bytes, err := json.Marshal(currentSettings)
	if err != nil {
		return err
	}

	err = c.db.Model(&dbModel).Updates(map[string]interface{}{"settings": bytes}).Error
	if err != nil {
		return err
	}

	return nil
}
