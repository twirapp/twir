package twirp_handlers

import (
	"github.com/satont/twir/apps/api-twirp/internal/interceptors"
	"github.com/satont/twir/apps/api-twirp/internal/wrappers"
	"github.com/satont/twir/libs/grpc/generated/api"
	"github.com/satont/twir/libs/grpc/generated/api/bots"
	"github.com/satont/twir/libs/grpc/generated/api/commands"
	"github.com/satont/twir/libs/grpc/generated/api/events"
	"github.com/satont/twir/libs/grpc/generated/api/greetings"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/twitchtv/twirp"
	"go.uber.org/fx"
	"time"
)

type Grpc struct {
	fx.In

	Tokens tokens.TokensClient
}

func NewProtected(opts Opts) IHandler {
	twirpHandler := api.NewProtectedServer(
		opts.ImplProtected,
		twirp.WithServerPathPrefix("/v1"),
		twirp.WithServerInterceptors(opts.Interceptor.DbUserInterceptor),
		twirp.WithServerInterceptors(opts.Interceptor.ChannelAccessInterceptor),
		twirp.WithServerInterceptors(opts.Interceptor.NewCacheInterceptor(
			interceptors.CacheOpts{
				CacheMethod:       "BotInfo",
				CacheDuration:     10 * time.Second,
				ClearMethods:      []string{"BotJoinPart"},
				WithChannelHeader: true,
				NewCastTo: func() any {
					return &bots.BotInfo{}
				},
			},
			interceptors.CacheOpts{
				CacheMethod:       "CommandsGetAll",
				CacheDuration:     24 * time.Hour,
				ClearMethods:      []string{"CommandsCreate", "CommandsDelete", "CommandsUpdate", "CommandsEnableOrDisable"},
				WithChannelHeader: true,
				NewCastTo: func() any {
					return &commands.CommandsGetAllResponse{}
				},
			},
			interceptors.CacheOpts{
				CacheMethod:       "EventsGetAll",
				CacheDuration:     24 * time.Hour,
				ClearMethods:      []string{"EventsCreate", "EventsDelete", "EventsUpdate", "EventsEnableOrDisable"},
				WithChannelHeader: true,
				NewCastTo: func() any {
					return &events.GetAllResponse{}
				},
			},
			interceptors.CacheOpts{
				CacheMethod:       "GreetingsGetAll",
				CacheDuration:     24 * time.Hour,
				ClearMethods:      []string{"GreetingsCreate", "GreetingsDelete", "GreetingsUpdate", "GreetingsEnableOrDisable"},
				WithChannelHeader: true,
				NewCastTo: func() any {
					return &greetings.GetAllResponse{}
				},
			},
		)),
		twirp.WithServerInterceptors(opts.Interceptor.Errors),
	)

	h := &Handler{
		pattern: twirpHandler.PathPrefix(),
		handler: wrappers.Wrap(
			twirpHandler,
			wrappers.WithCors,
			wrappers.WithDashboardId,
			wrappers.WithApiKeyHeader,
		),
	}

	return h
}
