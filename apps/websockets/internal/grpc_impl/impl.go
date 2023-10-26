package grpc_impl

import (
	"context"
	"fmt"
	"net"

	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/websockets/internal/namespaces/alerts"
	"github.com/satont/twir/apps/websockets/internal/namespaces/obs"
	"github.com/satont/twir/apps/websockets/internal/namespaces/registry/overlays"
	"github.com/satont/twir/apps/websockets/internal/namespaces/tts"
	"github.com/satont/twir/apps/websockets/internal/namespaces/youtube"
	"github.com/satont/twir/libs/grpc/constants"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"gorm.io/gorm"
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

	gorm   *gorm.DB
	redis  *redis.Client
	logger logger.Logger

	ttsServer              *tts.TTS
	youTubeServer          *youtube.YouTube
	obsServer              *obs.OBS
	alertsServer           *alerts.Alerts
	overlaysRegistryServer *overlays.Registry
}

type GrpcOpts struct {
	fx.In
	LC fx.Lifecycle

	Gorm   *gorm.DB
	Redis  *redis.Client
	Logger logger.Logger

	TTSServer              *tts.TTS
	YouTubeServer          *youtube.YouTube
	OBSServer              *obs.OBS
	AlertsServer           *alerts.Alerts
	OverlaysRegistryServer *overlays.Registry
}

func NewGrpcImplementation(opts GrpcOpts) error {
	impl := &GrpcImpl{
		gorm:                   opts.Gorm,
		redis:                  opts.Redis,
		logger:                 opts.Logger,
		ttsServer:              opts.TTSServer,
		youTubeServer:          opts.YouTubeServer,
		obsServer:              opts.OBSServer,
		alertsServer:           opts.AlertsServer,
		overlaysRegistryServer: opts.OverlaysRegistryServer,
	}

	grpcServer := grpc.NewServer()

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", constants.WEBSOCKET_SERVER_PORT))
				if err != nil {
					return err
				}
				websockets.RegisterWebsocketServer(grpcServer, impl)

				go grpcServer.Serve(lis)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				grpcServer.GracefulStop()
				return nil
			},
		},
	)

	return nil
}
