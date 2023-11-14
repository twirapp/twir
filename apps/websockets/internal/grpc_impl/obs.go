package grpc_impl

import (
	"context"

	"github.com/satont/twir/libs/grpc/generated/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *GrpcImpl) ObsSetScene(
	_ context.Context,
	msg *websockets.ObsSetSceneMessage,
) (*emptypb.Empty, error) {
	if err := c.obsServer.SendEvent(msg.ChannelId, "setScene", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *GrpcImpl) ObsToggleSource(
	_ context.Context,
	msg *websockets.ObsToggleSourceMessage,
) (*emptypb.Empty, error) {
	if err := c.obsServer.SendEvent(msg.ChannelId, "toggleSource", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *GrpcImpl) ObsToggleAudio(
	_ context.Context,
	msg *websockets.ObsToggleAudioMessage,
) (*emptypb.Empty, error) {
	if err := c.obsServer.SendEvent(msg.ChannelId, "toggleAudioSource", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *GrpcImpl) ObsAudioSetVolume(_ context.Context, msg *websockets.ObsAudioSetVolumeMessage) (
	*emptypb.Empty, error,
) {
	if err := c.obsServer.SendEvent(msg.ChannelId, "setVolume", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *GrpcImpl) ObsAudioIncreaseVolume(
	_ context.Context, msg *websockets.ObsAudioIncreaseVolumeMessage,
) (*emptypb.Empty, error) {
	if err := c.obsServer.SendEvent(msg.ChannelId, "increaseVolume", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *GrpcImpl) ObsAudioDecreaseVolume(
	_ context.Context, msg *websockets.ObsAudioDecreaseVolumeMessage,
) (*emptypb.Empty, error) {
	if err := c.obsServer.SendEvent(msg.ChannelId, "decreaseVolume", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *GrpcImpl) ObsAudioEnable(
	_ context.Context,
	msg *websockets.ObsAudioDisableOrEnableMessage,
) (
	*emptypb.Empty, error,
) {
	if err := c.obsServer.SendEvent(msg.ChannelId, "enableAudio", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *GrpcImpl) ObsAudioDisable(
	_ context.Context,
	msg *websockets.ObsAudioDisableOrEnableMessage,
) (
	*emptypb.Empty, error,
) {
	if err := c.obsServer.SendEvent(msg.ChannelId, "disableAudio", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *GrpcImpl) ObsStopStream(
	_ context.Context,
	msg *websockets.ObsStopOrStartStream,
) (*emptypb.Empty, error) {
	if err := c.obsServer.SendEvent(msg.ChannelId, "stopStream", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *GrpcImpl) ObsStartStream(
	_ context.Context,
	msg *websockets.ObsStopOrStartStream,
) (*emptypb.Empty, error) {
	if err := c.obsServer.SendEvent(msg.ChannelId, "startStream", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *GrpcImpl) ObsCheckIsUserConnected(
	_ context.Context,
	msg *websockets.ObsCheckUserConnectedRequest,
) (*websockets.ObsCheckUserConnectedResponse, error) {
	res, err := c.obsServer.IsUserConnected(msg.UserId)
	if err != nil {
		return nil, err
	}

	return &websockets.ObsCheckUserConnectedResponse{
		State: res,
	}, nil
}
