package grpc_impl

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/emotes-cacher/internal/emotes"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/emotes_cacher"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type EmotesCacherImpl struct {
	emotes_cacher.UnimplementedEmotesCacherServer

	redis  *redis.Client
	logger logger.Logger
}

type Opts struct {
	fx.In

	Redis  *redis.Client
	Logger logger.Logger
	Lc     fx.Lifecycle
}

func NewEmotesCacher(opts Opts) {
	impl := &EmotesCacherImpl{
		redis:  opts.Redis,
		logger: opts.Logger,
	}

	grpcServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				lis, err := net.Listen(
					"tcp",
					fmt.Sprintf("0.0.0.0:%d", constants.EMOTES_CACHER_SERVER_PORT),
				)
				if err != nil {
					return err
				}
				emotes_cacher.RegisterEmotesCacherServer(grpcServer, impl)
				go grpcServer.Serve(lis)

				return nil
			},
			OnStop: func(ctx context.Context) error {
				grpcServer.GracefulStop()
				return nil
			},
		},
	)
}

func (c *EmotesCacherImpl) CacheChannelEmotes(
	_ context.Context,
	req *emotes_cacher.Request,
) (*emptypb.Empty, error) {
	if req.ChannelId == "" {
		return &emptypb.Empty{}, nil
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	resultEmotes := make([]string, 0, 300)

	reqFuncs := []func(c string) ([]string, error){
		emotes.GetChannelSevenTvEmotes,
		emotes.GetChannelBttvEmotes,
		emotes.GetChannelFfzEmotes,
	}

	for _, f := range reqFuncs {
		wg.Add(1)
		f := f
		go func() {
			defer wg.Done()
			res, err := f(req.ChannelId)
			if err != nil {
				c.logger.Error("cannot get emotes", slog.Any("err", err))
				return
			}

			mu.Lock()
			resultEmotes = append(resultEmotes, res...)
			mu.Unlock()
		}()
	}

	wg.Wait()

	for _, emote := range resultEmotes {
		if emote == "" {
			continue
		}

		go func(emote string) {
			c.redis.Set(
				context.Background(),
				fmt.Sprintf("emotes:channel:%s:%s", req.ChannelId, emote),
				emote,
				10*time.Minute,
			)
		}(emote)
	}

	return &emptypb.Empty{}, nil
}

func (c *EmotesCacherImpl) CacheGlobalEmotes(_ context.Context, _ *emptypb.Empty) (
	*emptypb.Empty,
	error,
) {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	resultEmotes := make([]string, 300)

	wg.Add(3)

	go func() {
		defer wg.Done()
		em, err := emotes.GetGlobalSevenTvEmotes()
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		mu.Lock()
		defer mu.Unlock()
		resultEmotes = append(resultEmotes, em...)
	}()

	go func() {
		defer wg.Done()
		em, err := emotes.GetGlobalFfzEmotes()
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		mu.Lock()
		defer mu.Unlock()
		resultEmotes = append(resultEmotes, em...)
	}()

	go func() {
		defer wg.Done()
		em, err := emotes.GetGlobalBttvEmotes()
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		mu.Lock()
		defer mu.Unlock()
		resultEmotes = append(resultEmotes, em...)
	}()

	wg.Wait()

	for _, emote := range resultEmotes {
		if emote == "" {
			continue
		}
		go func(emote string) {
			c.redis.Set(
				context.Background(),
				fmt.Sprintf("emotes:global:%s", emote),
				emote,
				10*time.Minute,
			)
		}(emote)
	}

	return &emptypb.Empty{}, nil
}
