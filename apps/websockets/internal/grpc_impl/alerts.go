package grpc_impl

import (
	"context"
	"fmt"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *GrpcImpl) TriggerAlert(ctx context.Context, req *websockets.TriggerAlertRequest) (
	*emptypb.Empty, error,
) {
	entity := model.ChannelAlert{}
	if err := c.gorm.WithContext(ctx).Where(
		"channel_id = ? and id = ?", req.ChannelId,
		req.AlertId,
	).Find(&entity).Error; err != nil {
		return nil, err
	}

	if entity.ID == "" {
		return nil, fmt.Errorf(
			"cannot find alert with id %s and channel_id %s",
			req.AlertId,
			req.ChannelId,
		)
	}

	err := c.alertsServer.SendEvent(req.ChannelId, "trigger", entity)
	return &emptypb.Empty{}, err
}
