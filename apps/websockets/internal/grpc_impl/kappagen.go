package grpc_impl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *GrpcImpl) TriggerKappagen(
	_ context.Context,
	msg *websockets.TriggerKappagenRequest,
) (*emptypb.Empty, error) {
	if err := c.kappagenServer.SendEvent(msg.ChannelId, "kappagen", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *GrpcImpl) TriggerKappagenByEvent(
	ctx context.Context,
	req *websockets.TriggerKappagenByEventRequest,
) (*emptypb.Empty, error) {
	settings := &model.ChannelModulesSettings{}
	err := c.gorm.
		WithContext(ctx).
		Where(`"channelId" = ? AND "type" = ?`, req.ChannelId, "kappagen_overlay").
		Find(settings).
		Error
	if err != nil {
		return &emptypb.Empty{}, nil
	}

	if settings.ID == "" {
		return &emptypb.Empty{}, nil
	}

	parsedSettings := model.KappagenOverlaySettings{}
	if err := json.Unmarshal(settings.Settings, &parsedSettings); err != nil {
		return nil, fmt.Errorf("cannot parse kappagen settings %w", err)
	}

	ok := lo.SomeBy(
		parsedSettings.EnabledEvents, func(item int32) bool {
			return int32(req.Event.Number()) == item
		},
	)

	if ok {
		c.kappagenServer.SendEvent(req.ChannelId, "event", map[string]any{})
	}

	return &emptypb.Empty{}, nil
}
