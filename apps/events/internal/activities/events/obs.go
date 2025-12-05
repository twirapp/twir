package events

import (
	"context"
	"errors"
	"strconv"

	"github.com/twirapp/twir/apps/events/internal/shared"
	"github.com/twirapp/twir/libs/bus-core/api"
	"github.com/twirapp/twir/libs/repositories/events/model"
	"go.temporal.io/sdk/activity"
)

func (c *Activity) ObsSetScene(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)
	if operation.Input == nil || *operation.Input == "" {
		return errors.New("input is required for operation ObsSetScene")
	}
	hydratedString, hydratedErr := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		*operation.Input,
		data,
	)
	if hydratedErr != nil {
		return hydratedErr
	}
	if hydratedString == "" {
		return nil
	}
	err := c.bus.Api.TriggerObsCommand.Publish(
		ctx,
		api.TriggerObsCommand{
			ChannelId: data.ChannelID,
			Action:    api.ObsCommandActionSetScene,
			Target:    hydratedString,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
func (c *Activity) ObsToggleSource(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)
	if operation.Target == nil || *operation.Target == "" {
		return errors.New("target is required for operation ObsToggleSource")
	}
	hydratedString, hydratedErr := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		*operation.Target,
		data,
	)
	if hydratedErr != nil {
		return hydratedErr
	}
	if hydratedString == "" {
		return nil
	}
	err := c.bus.Api.TriggerObsCommand.Publish(
		ctx,
		api.TriggerObsCommand{
			ChannelId: data.ChannelID,
			Action:    api.ObsCommandActionToggleSource,
			Target:    hydratedString,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
func (c *Activity) ObsToggleAudio(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	if operation.Target == nil || *operation.Target == "" {
		return errors.New("target is required for operation ObsToggleAudio")
	}
	hydratedString, hydratedErr := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		*operation.Target,
		data,
	)
	if hydratedErr != nil {
		return hydratedErr
	}
	if hydratedString == "" {
		return nil
	}
	err := c.bus.Api.TriggerObsCommand.Publish(
		ctx,
		api.TriggerObsCommand{
			ChannelId: data.ChannelID,
			Action:    api.ObsCommandActionToggleAudio,
			Target:    hydratedString,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
func (c *Activity) ObsAudioChangeVolume(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)
	if operation.Target == nil || *operation.Target == "" {
		return errors.New("target is required for operation ObsAudioChangeVolume")
	}
	msg, err := c.hydrator.HydrateStringWithData(data.ChannelID, *operation.Target, data)
	if err != nil {
		return err
	}
	if msg == "" {
		return nil
	}
	parsedStep, err := strconv.Atoi(msg)
	if err != nil {
		return err
	}
	if parsedStep < 0 || parsedStep > 20 {
		return errors.New("step must be between 0 and 20")
	}
	var action api.ObsCommandAction
	if operation.Type == model.EventOperationTypeObsIncreaseAudioVolume {
		action = api.ObsCommandActionIncreaseVolume
	} else {
		action = api.ObsCommandActionDecreaseVolume
	}
	err = c.bus.Api.TriggerObsCommand.Publish(
		ctx,
		api.TriggerObsCommand{
			ChannelId:  data.ChannelID,
			Action:     action,
			Target:     *operation.Target,
			VolumeStep: &parsedStep,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
func (c *Activity) ObsAudioSetVolume(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)
	if operation.Target == nil || *operation.Target == "" {
		return errors.New("target is required for operation ObsAudioSetVolume")
	}
	msg, err := c.hydrator.HydrateStringWithData(data.ChannelID, *operation.Target, data)
	if err != nil {
		return err
	}
	if msg == "" {
		return nil
	}
	parsedVolume, err := strconv.Atoi(msg)
	if err != nil {
		return err
	}
	if parsedVolume < 0 || parsedVolume > 20 {
		return errors.New("volume must be between 0 and 20")
	}
	err = c.bus.Api.TriggerObsCommand.Publish(
		ctx,
		api.TriggerObsCommand{
			ChannelId:   data.ChannelID,
			Action:      api.ObsCommandActionSetVolume,
			Target:      *operation.Target,
			VolumeValue: &parsedVolume,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
func (c *Activity) ObsEnableOrDisableAudio(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)
	if operation.Target == nil || *operation.Target == "" {
		return errors.New("target is required for operation ObsEnableOrDisableAudio")
	}
	var action api.ObsCommandAction
	if operation.Type == model.EventOperationTypeObsDisableAudio {
		action = api.ObsCommandActionDisableAudio
	} else {
		action = api.ObsCommandActionEnableAudio
	}
	err := c.bus.Api.TriggerObsCommand.Publish(
		ctx,
		api.TriggerObsCommand{
			ChannelId: data.ChannelID,
			Action:    action,
			Target:    *operation.Target,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
func (c *Activity) ObsStartOrStopStream(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)
	var action api.ObsCommandAction
	if operation.Type == model.EventOperationTypeObsStartStream {
		action = api.ObsCommandActionStartStream
	} else {
		action = api.ObsCommandActionStopStream
	}
	err := c.bus.Api.TriggerObsCommand.Publish(
		ctx,
		api.TriggerObsCommand{
			ChannelId: data.ChannelID,
			Action:    action,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
