package impl

import (
	"context"
	"github.com/satont/tsuwari/libs/grpc/generated/api/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (a Api) BotInfo(ctx context.Context, meta *api.BaseRequestMeta) (*api.BotInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (a Api) BotJoinPart(ctx context.Context, request *api.BotJoinPartRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
