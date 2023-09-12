package grpc_impl

import (
	"github.com/satont/twir/apps/websockets/internal/namespaces/alerts"
	"github.com/satont/twir/apps/websockets/internal/namespaces/obs"
	"github.com/satont/twir/apps/websockets/internal/namespaces/registry/overlays"
	"github.com/satont/twir/apps/websockets/internal/namespaces/tts"
	"github.com/satont/twir/apps/websockets/internal/namespaces/youtube"
	"github.com/satont/twir/apps/websockets/types"
	"github.com/satont/twir/libs/grpc/generated/websockets"
)

type Sockets struct {
	TTS              *tts.TTS
	YouTube          *youtube.YouTube
	OBS              *obs.OBS
	Alerts           *alerts.Alerts
	OverlaysRegistry *overlays.Registry
}

type GrpcImpl struct {
	websockets.UnimplementedWebsocketServer
	sockets  *Sockets
	services *types.Services
}

type GrpcOpts struct {
	Services *types.Services
	Sockets  *Sockets
}

func NewGrpcImplementation(opts *GrpcOpts) *GrpcImpl {
	return &GrpcImpl{
		sockets:  opts.Sockets,
		services: opts.Services,
	}
}
