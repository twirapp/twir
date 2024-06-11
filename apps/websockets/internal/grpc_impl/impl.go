package grpc_impl

import (
	"context"
	"fmt"
	"net"

	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/websockets/internal/namespaces/overlays/alerts"
	"github.com/satont/twir/apps/websockets/internal/namespaces/overlays/be_right_back"
	"github.com/satont/twir/apps/websockets/internal/namespaces/overlays/dudes"
	"github.com/satont/twir/apps/websockets/internal/namespaces/overlays/kappagen"
	"github.com/satont/twir/apps/websockets/internal/namespaces/overlays/obs"
	"github.com/satont/twir/apps/websockets/internal/namespaces/overlays/registry/overlays"
	"github.com/satont/twir/apps/websockets/internal/namespaces/overlays/tts"
	"github.com/satont/twir/apps/websockets/internal/namespaces/youtube"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
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
	kappagenServer         *kappagen.Kappagen
	beRightBackServer      *be_right_back.BeRightBack
	dudesServer            *dudes.Dudes
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
	KappagenServer         *kappagen.Kappagen
	BeRightBackServer      *be_right_back.BeRightBack
	DudesServer            *dudes.Dudes
}

func NewGrpcImplementation(opts GrpcOpts) (websockets.WebsocketServer, error) {
	impl := &GrpcImpl{
		gorm:                   opts.Gorm,
		redis:                  opts.Redis,
		logger:                 opts.Logger,
		ttsServer:              opts.TTSServer,
		youTubeServer:          opts.YouTubeServer,
		obsServer:              opts.OBSServer,
		alertsServer:           opts.AlertsServer,
		overlaysRegistryServer: opts.OverlaysRegistryServer,
		kappagenServer:         opts.KappagenServer,
		beRightBackServer:      opts.BeRightBackServer,
		dudesServer:            opts.DudesServer,
	}

	grpcServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))

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

	return impl, nil
}

func (c *GrpcImpl) RefreshOverlaySettings(
	_ context.Context,
	req *websockets.RefreshOverlaysRequest,
) (
	*emptypb.Empty,
	error,
) {
	var err error

	switch req.GetOverlayName() {
	case websockets.RefreshOverlaySettingsName_CUSTOM:
		err = c.overlaysRegistryServer.SendEvent(
			req.GetChannelId(),
			"refreshOverlays",
			nil,
		)
	case websockets.RefreshOverlaySettingsName_KAPPAGEN:
		err = c.kappagenServer.SendSettings(req.GetChannelId())
	case websockets.RefreshOverlaySettingsName_BRB:
		err = c.beRightBackServer.SendSettings(req.GetChannelId())
	case websockets.RefreshOverlaySettingsName_DUDES:
		err = c.dudesServer.SendSettings(req.GetChannelId(), req.GetOverlayId())
	default:
		return nil, fmt.Errorf("unknown overlay: %s", req.GetOverlayName())
	}

	if err != nil {
		c.logger.Error(err.Error())
	}

	return &emptypb.Empty{}, nil
}
