package resolvers

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/minio/minio-go/v7"
	"github.com/nicklaw5/helix/v2"
	"github.com/redis/go-redis/v9"
	config "github.com/satont/twir/libs/config"
	deprecatedgormmodel "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	twir_stats "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/twir-stats"
	admin_actions "github.com/twirapp/twir/apps/api-gql/internal/services/admin-actions"
	"github.com/twirapp/twir/apps/api-gql/internal/services/alerts"
	audit_logs "github.com/twirapp/twir/apps/api-gql/internal/services/audit-logs"
	"github.com/twirapp/twir/apps/api-gql/internal/services/badges"
	badges_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-users"
	dashboard_widget_events "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard-widget-events"
	"github.com/twirapp/twir/apps/api-gql/internal/services/keywords"
	"github.com/twirapp/twir/apps/api-gql/internal/services/timers"
	twir_users "github.com/twirapp/twir/apps/api-gql/internal/services/twir-users"
	"github.com/twirapp/twir/apps/api-gql/internal/services/users"
	"github.com/twirapp/twir/apps/api-gql/internal/services/variables"
	"github.com/twirapp/twir/apps/api-gql/internal/wsrouter"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	twitchcahe "github.com/twirapp/twir/libs/cache/twitch"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	config               config.Config
	sessions             *auth.Auth
	gorm                 *gorm.DB
	twitchClient         *helix.Client
	cachedTwitchClient   *twitchcahe.CachedTwitchClient
	cachedCommandsClient *generic_cacher.GenericCacher[[]deprecatedgormmodel.ChannelsCommands]
	minioClient          *minio.Client
	twirBus              *bus_core.Bus
	logger               logger.Logger
	redis                *redis.Client
	tokensClient         tokens.TokensClient
	wsRouter             wsrouter.WsRouter
	twirStats            *twir_stats.TwirStats

	dashboardWidgetEventsService *dashboard_widget_events.Service
	variablesService             *variables.Service
	timersService                *timers.Service
	keywordsService              *keywords.Service
	auditLogService              *audit_logs.Service
	adminActionsService          *admin_actions.Service
	badgesService                *badges.Service
	badgesUsersService           *badges_users.Service
	usersService                 *users.Service
	twirUsersService             *twir_users.Service
	alertsService                *alerts.Service
}

type Opts struct {
	fx.In

	Sessions             *auth.Auth
	Gorm                 *gorm.DB
	Config               config.Config
	TokensGrpc           tokens.TokensClient
	CachedTwitchClient   *twitchcahe.CachedTwitchClient
	CachedCommandsClient *generic_cacher.GenericCacher[[]deprecatedgormmodel.ChannelsCommands]
	Minio                *minio.Client
	TwirBus              *bus_core.Bus
	Logger               logger.Logger
	Redis                *redis.Client
	WsRouter             wsrouter.WsRouter
	TwirStats            *twir_stats.TwirStats

	DashboardWidgetEventsService *dashboard_widget_events.Service
	VariablesService             *variables.Service
	TimersService                *timers.Service
	KeywordService               *keywords.Service
	UserAuditLogService          *audit_logs.Service
	AdminActionsService          *admin_actions.Service
	BadgesService                *badges.Service
	BadgesUsersService           *badges_users.Service
	UsersService                 *users.Service
	TwirUsersService             *twir_users.Service
	AlertsService                *alerts.Service
}

func New(opts Opts) (*Resolver, error) {
	twitchClient, err := twitch.NewAppClient(opts.Config, opts.TokensGrpc)
	if err != nil {
		return nil, err
	}

	return &Resolver{
		config:                       opts.Config,
		sessions:                     opts.Sessions,
		gorm:                         opts.Gorm,
		twitchClient:                 twitchClient,
		cachedTwitchClient:           opts.CachedTwitchClient,
		minioClient:                  opts.Minio,
		twirBus:                      opts.TwirBus,
		logger:                       opts.Logger,
		redis:                        opts.Redis,
		cachedCommandsClient:         opts.CachedCommandsClient,
		tokensClient:                 opts.TokensGrpc,
		wsRouter:                     opts.WsRouter,
		twirStats:                    opts.TwirStats,
		dashboardWidgetEventsService: opts.DashboardWidgetEventsService,
		variablesService:             opts.VariablesService,
		timersService:                opts.TimersService,
		keywordsService:              opts.KeywordService,
		auditLogService:              opts.UserAuditLogService,
		adminActionsService:          opts.AdminActionsService,
		badgesService:                opts.BadgesService,
		badgesUsersService:           opts.BadgesUsersService,
		usersService:                 opts.UsersService,
		twirUsersService:             opts.TwirUsersService,
		alertsService:                opts.AlertsService,
	}, nil
}

func GetPreloads(ctx context.Context) []string {
	return GetNestedPreloads(
		graphql.GetOperationContext(ctx),
		graphql.CollectFieldsCtx(ctx, nil),
		"",
	)
}

func GetNestedPreloads(
	ctx *graphql.OperationContext,
	fields []graphql.CollectedField,
	prefix string,
) (preloads []string) {
	for _, column := range fields {
		prefixColumn := GetPreloadString(prefix, column.Name)
		preloads = append(preloads, prefixColumn)
		preloads = append(
			preloads,
			GetNestedPreloads(ctx, graphql.CollectFields(ctx, column.Selections, nil), prefixColumn)...,
		)
	}
	return
}

func GetPreloadString(prefix, name string) string {
	if len(prefix) > 0 {
		return prefix + "." + name
	}
	return name
}
