package events

import (
	"context"
	"errors"
	"strconv"

	"github.com/twirapp/twir/apps/events/internal/shared"
	"github.com/twirapp/twir/libs/grpc/websockets"
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

	_, err := c.websocketsGrpc.ObsSetScene(
		ctx,
		&websockets.ObsSetSceneMessage{
			ChannelId: data.ChannelID,
			SceneName: hydratedString,
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

	_, err := c.websocketsGrpc.ObsToggleSource(
		ctx,
		&websockets.ObsToggleSourceMessage{
			ChannelId:  data.ChannelID,
			SourceName: hydratedString,
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

	_, err := c.websocketsGrpc.ObsToggleAudio(
		context.Background(), &websockets.ObsToggleAudioMessage{
			ChannelId:       data.ChannelID,
			AudioSourceName: hydratedString,
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

	if operation.Type == model.EventOperationTypeObsIncreaseAudioVolume {
		_, err = c.websocketsGrpc.ObsAudioIncreaseVolume(
			ctx,
			&websockets.ObsAudioIncreaseVolumeMessage{
				ChannelId:       data.ChannelID,
				AudioSourceName: *operation.Target,
				Step:            uint32(parsedStep),
			},
		)
		if err != nil {
			return err
		}
	} else {
		_, err = c.websocketsGrpc.ObsAudioDecreaseVolume(
			context.Background(), &websockets.ObsAudioDecreaseVolumeMessage{
				ChannelId:       data.ChannelID,
				AudioSourceName: *operation.Target,
				Step:            uint32(parsedStep),
			},
		)
		if err != nil {
			return err
		}
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

	_, err = c.websocketsGrpc.ObsAudioSetVolume(
		ctx,
		&websockets.ObsAudioSetVolumeMessage{
			ChannelId:       data.ChannelID,
			AudioSourceName: *operation.Target,
			Volume:          uint32(parsedVolume),
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

	if operation.Type == model.EventOperationTypeObsDisableAudio {
		_, err := c.websocketsGrpc.ObsAudioDisable(
			ctx,
			&websockets.ObsAudioDisableOrEnableMessage{
				ChannelId:       data.ChannelID,
				AudioSourceName: *operation.Target,
			},
		)
		if err != nil {
			return err
		}
	}

	if operation.Type == model.EventOperationTypeObsEnableAudio {
		_, err := c.websocketsGrpc.ObsAudioEnable(
			ctx,
			&websockets.ObsAudioDisableOrEnableMessage{
				ChannelId:       data.ChannelID,
				AudioSourceName: *operation.Target,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Activity) ObsStartOrStopStream(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	if operation.Type == model.EventOperationTypeObsStartStream {
		_, err := c.websocketsGrpc.ObsStartStream(
			ctx,
			&websockets.ObsStopOrStartStream{
				ChannelId: data.ChannelID,
			},
		)
		if err != nil {
			return err
		}
	}

	if operation.Type == model.EventOperationTypeObsStopStream {
		_, err := c.websocketsGrpc.ObsStopStream(
			ctx,
			&websockets.ObsStopOrStartStream{
				ChannelId: data.ChannelID,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
