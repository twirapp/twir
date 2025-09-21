package services

import (
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/go-redsync/redsync/v4"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/parser/internal/services/chat_wall"
	"github.com/twirapp/twir/apps/parser/internal/services/shortenedurls"
	ttsservice "github.com/twirapp/twir/apps/parser/internal/services/tts"
	"github.com/twirapp/twir/apps/parser/pkg/executron"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/cache/twitch"
	cfg "github.com/twirapp/twir/libs/config"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/i18n"
	seventvintegrationapi "github.com/twirapp/twir/libs/integrations/seventv/api"
	channelscategoriesaliases "github.com/twirapp/twir/libs/repositories/channels_categories_aliases"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	channelscommandsusages "github.com/twirapp/twir/libs/repositories/channels_commands_usages"
	channelsemotesusagesrepository "github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	channelsinfohistory "github.com/twirapp/twir/libs/repositories/channels_info_history"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	channelsmodules_settingstts "github.com/twirapp/twir/libs/repositories/channels_modules_settings_tts"
	chatmessagesrepository "github.com/twirapp/twir/libs/repositories/chat_messages"
	chatwallrepository "github.com/twirapp/twir/libs/repositories/chat_wall"
	chatwallmodel "github.com/twirapp/twir/libs/repositories/chat_wall/model"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	"github.com/twirapp/twir/libs/repositories/users"
	"github.com/twirapp/twir/libs/types/types/api/modules"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Grpc struct {
	WebSockets websockets.WebsocketClient
}

type Services struct {
	Config                     *cfg.Config
	Logger                     *zap.Logger
	Gorm                       *gorm.DB
	Sqlx                       *sqlx.DB
	Redis                      *redis.Client
	GrpcClients                *Grpc
	Bus                        *buscore.Bus
	TrmManager                 trm.Manager
	CommandsCache              *generic_cacher.GenericCacher[[]model.ChannelsCommands]
	CommandsPrefixCache        *generic_cacher.GenericCacher[channelscommandsprefixmodel.ChannelsCommandsPrefix]
	SevenTvCache               *generic_cacher.GenericCacher[seventvintegrationapi.TwirSeventvUser]
	ChatWallCache              *generic_cacher.GenericCacher[[]chatwallmodel.ChatWall]
	ChatWallService            *chat_wall.Service
	RedSync                    *redsync.Redsync
	CommandsLock               *redsync.Mutex
	CommandsPrefixRepository   channelscommandsprefixrepository.Repository
	TTSCache                   *generic_cacher.GenericCacher[modules.TTSSettings]
	TTSRepository              channelsmodules_settingstts.Repository
	TTSService                 *ttsservice.Service
	SpotifyRepo                channelsintegrationsspotify.Repository
	UsersRepo                  users.Repository
	CategoriesAliasesRepo      channelscategoriesaliases.Repository
	ScheduledVipsRepo          scheduledvipsrepository.Repository
	CacheTwitchClient          *twitch.CachedTwitchClient
	ChatWallRepo               chatwallrepository.Repository
	ChannelsInfoHistoryRepo    channelsinfohistory.Repository
	ChannelEmotesUsagesRepo    channelsemotesusagesrepository.Repository
	ChannelsCommandsUsagesRepo channelscommandsusages.Repository
	ChatMessagesRepo           chatmessagesrepository.Repository
	ChannelEventListsRepo      channelseventslist.Repository
	ShortUrlServices           *shortenedurls.Service
	Executron                  executron.Executron
	I18n                       *i18n.I18n
}
