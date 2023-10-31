package twirp_handlers

import (
	"time"

	"github.com/satont/twir/apps/api/internal/handlers"
	"github.com/satont/twir/apps/api/internal/hooks"
	"github.com/satont/twir/apps/api/internal/interceptors"
	"github.com/satont/twir/apps/api/internal/wrappers"
	"github.com/satont/twir/libs/grpc/generated/api"
	"github.com/satont/twir/libs/grpc/generated/api/google_fonts_unprotected"
	"github.com/twitchtv/twirp"
)

func NewUnProtected(opts Opts) handlers.IHandler {
	twirpHandler := api.NewUnProtectedServer(
		opts.ImplUnProtected,
		twirp.WithServerPathPrefix("/v1"),
		twirp.WithServerHooks(hooks.NewLoggingServerHooks(opts.Logger)),
		twirp.WithServerInterceptors(
			opts.Interceptor.NewCacheInterceptor(
				interceptors.CacheOpts{
					CacheMethod:   "GetGoogleFonts",
					CacheDuration: 12 * time.Hour,
					CastTo: func() any {
						return &google_fonts_unprotected.Fonts{}
					},
				},
				// interceptors.CacheOpts{
				//	CacheMethod:       "CommandsGetAll",
				//	CacheDuration:     24 * time.Hour,
				//	ClearMethods:      []string{"CommandsCreate", "CommandsDelete", "CommandsUpdate", "CommandsEnableOrDisable"},
				//	WithChannelHeader: true,
				//	CastTo: func() any {
				//		return &commands.CommandsGetAllResponse{}
				//	},
				// },
				// interceptors.CacheOpts{
				//	CacheMethod:       "EventsGetAll",
				//	CacheDuration:     24 * time.Hour,
				//	ClearMethods:      []string{"EventsCreate", "EventsDelete", "EventsUpdate", "EventsEnableOrDisable"},
				//	WithChannelHeader: true,
				//	CastTo: func() any {
				//		return &events.GetAllResponse{}
				//	},
				// },
				// interceptors.CacheOpts{
				//	CacheMethod:       "GreetingsGetAll",
				//	CacheDuration:     24 * time.Hour,
				//	ClearMethods:      []string{"GreetingsCreate", "GreetingsDelete", "GreetingsUpdate", "GreetingsEnableOrDisable"},
				//	WithChannelHeader: true,
				//	CastTo: func() any {
				//		return &greetings.GetAllResponse{}
				//	},
				// },
				// interceptors.CacheOpts{
				//	CacheMethod:       "KeywordsGetAll",
				//	CacheDuration:     24 * time.Hour,
				//	ClearMethods:      []string{"KeywordsCreate", "KeywordsDelete", "KeywordsUpdate", "KeywordsEnableOrDisable"},
				//	WithChannelHeader: true,
				//	CastTo: func() any {
				//		return &keywords.GetAllResponse{}
				//	},
				// },
				// interceptors.CacheOpts{
				//	CacheMethod:       "TimersGetAll",
				//	CacheDuration:     24 * time.Hour,
				//	ClearMethods:      []string{"TimersCreate", "TimersDelete", "TimersUpdate", "TimersEnableOrDisable"},
				//	WithChannelHeader: true,
				//	CastTo: func() any {
				//		return &timers.GetResponse{}
				//	},
				// },
				// interceptors.CacheOpts{
				// 	CacheMethod:       "VariablesGetAll",
				// 	CacheDuration:     24 * time.Hour,
				// 	ClearMethods:      []string{"VariablesCreate", "VariablesDelete", "VariablesUpdate", "VariablesEnableOrDisable"},
				// 	WithChannelHeader: true,
				// 	NewCastTo: func() any {
				// 		return &variables.GetAllResponse{}
				// 	},
				// },
				// interceptors.CacheOpts{
				// 	CacheMethod:       "CommandsGroupGetAll",
				// 	CacheDuration:     24 * time.Hour,
				// 	ClearMethods:      []string{"CommandsGroupCreate", "CommandsGroupDelete", "CommandsGroupUpdate", "CommandsGroupEnableOrDisable"},
				// 	WithChannelHeader: true,
				// 	NewCastTo: func() any {
				// 		return &commands_group.GetAllResponse{}
				// 	},
				// },
				// interceptors.CacheOpts{
				// 	CacheMethod:       "RolesGetAll",
				// 	CacheDuration:     24 * time.Hour,
				// 	ClearMethods:      []string{"RolesGroupCreate", "RolesGroupDelete", "RolesGroupUpdate", "RolesGroupEnableOrDisable"},
				// 	WithChannelHeader: true,
				// 	NewCastTo: func() any {
				// 		return &roles.GetAllResponse{}
				// 	},
				// },
			),
		),
	)

	h := handlers.New(
		handlers.Opts{
			Pattern: twirpHandler.PathPrefix(),
			Handler: wrappers.Wrap(
				twirpHandler,
				wrappers.WithCors,
				wrappers.WithApiKeyHeader,
			),
		},
	)

	return h
}
