package grpc_impl

import (
	"github.com/satont/tsuwari/libs/grpc/generated/giveaways"
)

type giveawaysGrpcServer struct {
	giveaways.UnimplementedGrpcServer
}
