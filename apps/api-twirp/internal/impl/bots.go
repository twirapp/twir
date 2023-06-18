package impl

import (
	"context"
	"fmt"
	"github.com/satont/tsuwari/libs/grpc/generated/api/api"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func (a *Api) BotInfo(ctx context.Context, meta *api.BaseRequestMeta) (*api.BotInfo, error) {
	return &api.BotInfo{
		IsMod:   false,
		BotId:   fmt.Sprintf("%v", time.Now().Unix()),
		BotName: "",
		Enabled: false,
	}, nil
}

func (a *Api) BotJoinPart(ctx context.Context, request *api.BotJoinPartRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
