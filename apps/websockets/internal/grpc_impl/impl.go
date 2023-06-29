package grpc_impl

import (
	"github.com/satont/twir/apps/websockets/internal/namespaces"
	"github.com/satont/twir/apps/websockets/internal/namespaces/obs"
	"github.com/satont/twir/apps/websockets/internal/namespaces/youtube"
	"github.com/satont/twir/apps/websockets/types"
	"github.com/satont/twir/libs/grpc/generated/websockets"
)

type Sockets struct {
	TTS     *namespaces.NameSpace
	YouTube *youtube.YouTube
	OBS     *obs.OBS
}

type grpcImpl struct {
	websockets.UnimplementedWebsocketServer
	sockets  *Sockets
	services *types.Services
}

type GrpcOpts struct {
	Services *types.Services
	Sockets  *Sockets
}

func NewGrpcImplementation(opts *GrpcOpts) *grpcImpl {
	return &grpcImpl{
		sockets:  opts.Sockets,
		services: opts.Services,
	}
}
