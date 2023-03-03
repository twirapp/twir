package grpc_impl

import (
	"context"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *grpcImpl) ObsSetScene(_ context.Context, msg *websockets.ObsSetSceneMessage) (*emptypb.Empty, error) {
	if err := c.sockets.OBS.SendEvent(msg.ChannelId, "setScene", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *grpcImpl) ObsToggleSource(_ context.Context, msg *websockets.ObsToggleSourceMessage) (*emptypb.Empty, error) {
	if err := c.sockets.OBS.SendEvent(msg.ChannelId, "toggleSource", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *grpcImpl) ObsToggleAudio(_ context.Context, msg *websockets.ObsToggleAudioMessage) (*emptypb.Empty, error) {
	if err := c.sockets.OBS.SendEvent(msg.ChannelId, "toggleAudioSource", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *grpcImpl) ObsAudioSetVolume(_ context.Context, msg *websockets.ObsAudioSetVolumeMessage) (*emptypb.Empty, error) {
	if err := c.sockets.OBS.SendEvent(msg.ChannelId, "setVolume", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *grpcImpl) ObsAudioIncreaseVolume(_ context.Context, msg *websockets.ObsAudioIncreaseVolumeMessage) (*emptypb.Empty, error) {
	if err := c.sockets.OBS.SendEvent(msg.ChannelId, "increaseVolume", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *grpcImpl) ObsAudioDecreaseVolume(_ context.Context, msg *websockets.ObsAudioDecreaseVolumeMessage) (*emptypb.Empty, error) {
	if err := c.sockets.OBS.SendEvent(msg.ChannelId, "decreaseVolume", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *grpcImpl) ObsAudioEnable(_ context.Context, msg *websockets.ObsAudioDisableOrEnableMessage) (*emptypb.Empty, error) {
	if err := c.sockets.OBS.SendEvent(msg.ChannelId, "enableAudio", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *grpcImpl) ObsAudioDisable(_ context.Context, msg *websockets.ObsAudioDisableOrEnableMessage) (*emptypb.Empty, error) {
	if err := c.sockets.OBS.SendEvent(msg.ChannelId, "disableAudio", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *grpcImpl) ObsStopStream(_ context.Context, msg *websockets.ObsStopOrStartStream) (*emptypb.Empty, error) {
	if err := c.sockets.OBS.SendEvent(msg.ChannelId, "stopStream", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *grpcImpl) ObsStartStream(_ context.Context, msg *websockets.ObsStopOrStartStream) (*emptypb.Empty, error) {
	if err := c.sockets.OBS.SendEvent(msg.ChannelId, "startStream", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
