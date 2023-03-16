package processor

import (
	"context"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"
	"strconv"
)

func (c *Processor) ObsSetScene(input string) error {
	_, err := c.services.WebsocketsGrpc.ObsSetScene(context.Background(), &websockets.ObsSetSceneMessage{
		ChannelId: c.channelId,
		SceneName: input,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Processor) ObsToggleSource(input string) error {
	_, err := c.services.WebsocketsGrpc.ObsToggleSource(context.Background(), &websockets.ObsToggleSourceMessage{
		ChannelId:  c.channelId,
		SourceName: input,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Processor) ObsToggleAudio(sourceName string) error {
	_, err := c.services.WebsocketsGrpc.ObsToggleAudio(context.Background(), &websockets.ObsToggleAudioMessage{
		ChannelId:       c.channelId,
		AudioSourceName: sourceName,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Processor) ObsAudioChangeVolume(operationType model.EventOperationType, sourceName, input string) error {
	msg, err := c.hydrateStringWithData(input, c.data)
	if err != nil {
		return err
	}

	parsedStep, err := strconv.Atoi(msg)
	if err != nil {
		return err
	}

	if parsedStep < 0 || parsedStep > 20 {
		return InternalError
	}

	if operationType == model.OperationObsIncreaseVolume {
		_, err = c.services.WebsocketsGrpc.ObsAudioIncreaseVolume(context.Background(), &websockets.ObsAudioIncreaseVolumeMessage{
			ChannelId:       c.channelId,
			AudioSourceName: sourceName,
			Step:            uint32(parsedStep),
		})
		if err != nil {
			return err
		}
	} else {
		_, err = c.services.WebsocketsGrpc.ObsAudioDecreaseVolume(context.Background(), &websockets.ObsAudioDecreaseVolumeMessage{
			ChannelId:       c.channelId,
			AudioSourceName: sourceName,
			Step:            uint32(parsedStep),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Processor) ObsAudioSetVolume(sourceName, input string) error {
	msg, err := c.hydrateStringWithData(input, c.data)
	if err != nil {
		return err
	}

	parsedVolume, err := strconv.Atoi(msg)
	if err != nil {
		return err
	}

	if parsedVolume < 0 || parsedVolume > 20 {
		return InternalError
	}

	_, err = c.services.WebsocketsGrpc.ObsAudioSetVolume(context.Background(), &websockets.ObsAudioSetVolumeMessage{
		ChannelId:       c.channelId,
		AudioSourceName: sourceName,
		Volume:          uint32(parsedVolume),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Processor) ObsEnableOrDisableAudio(operation model.EventOperationType, sourceName string) error {
	if operation == model.OperationObsDisableAudio {
		_, err := c.services.WebsocketsGrpc.ObsAudioDisable(context.Background(), &websockets.ObsAudioDisableOrEnableMessage{
			ChannelId:       c.channelId,
			AudioSourceName: sourceName,
		})
		if err != nil {
			return err
		}
	}

	if operation == model.OperationObsEnableAudio {
		_, err := c.services.WebsocketsGrpc.ObsAudioEnable(context.Background(), &websockets.ObsAudioDisableOrEnableMessage{
			ChannelId:       c.channelId,
			AudioSourceName: sourceName,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Processor) ObsStartOrStopStream(operation model.EventOperationType) error {
	if operation == model.OperationObsStartStream {
		_, err := c.services.WebsocketsGrpc.ObsStartStream(context.Background(), &websockets.ObsStopOrStartStream{
			ChannelId: c.channelId,
		})
		if err != nil {
			return err
		}
	}

	if operation == model.OperationObsStopStream {
		_, err := c.services.WebsocketsGrpc.ObsStopStream(context.Background(), &websockets.ObsStopOrStartStream{
			ChannelId: c.channelId,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
