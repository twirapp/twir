package grpc_impl

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/grpc/websockets"
	alertmodel "github.com/twirapp/twir/libs/repositories/alerts/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *GrpcImpl) TriggerAlert(ctx context.Context, req *websockets.TriggerAlertRequest) (
	*emptypb.Empty, error,
) {
	alerts, err := c.alertsCache.Get(ctx, req.ChannelId)
	if err != nil {
		return nil, err
	}

	var foundAlert alertmodel.Alert
	for _, alert := range alerts {
		if alert.ID.String() == req.AlertId {
			foundAlert = alert
			break
		}
	}
	if foundAlert.ID == uuid.Nil {
		return nil, fmt.Errorf("cannot find alert with id %s", req.AlertId)
	}

	err = c.alertsServer.SendEvent(req.ChannelId, "trigger", foundAlert)
	return &emptypb.Empty{}, err
}
