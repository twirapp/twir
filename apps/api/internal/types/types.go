package types

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	rdb "github.com/go-redis/redis/v9"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"github.com/satont/tsuwari/apps/api/internal/services/redis"
	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"
	"github.com/satont/tsuwari/libs/grpc/generated/eventsub"
	"github.com/satont/tsuwari/libs/grpc/generated/integrations"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
	"github.com/satont/tsuwari/libs/grpc/generated/scheduler"
	"github.com/satont/tsuwari/libs/grpc/generated/timers"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"gorm.io/gorm"
)

type GrpcClientsService struct {
	Integrations integrations.IntegrationsClient
	Parser       parser.ParserClient
	EventSub     eventsub.EventSubClient
	Scheduler    scheduler.SchedulerClient
	Timers       timers.TimersClient
	Bots         bots.BotsClient
	Tokens       tokens.TokensClient
}

type Services struct {
	Logger              interfaces.Logger
	Redis               *rdb.Client
	Gorm                *gorm.DB
	Sqlx                *sqlx.DB
	Config              *cfg.Config
	RedisStorage        *redis.RedisStorage
	Validator           *validator.Validate
	ValidatorTranslator ut.Translator
	TgBotApi            *tgbotapi.BotAPI
	Grpc                *GrpcClientsService

	TimersService interfaces.TimersService
}

type JSONResult struct{}

type DOCApiBadRequest struct {
	Messages string
}

type DOCApiValidationError struct {
	Messages []string
}

type DOCApiNotFoundError struct {
	Messages []string
}

type DOCApiInternalError struct {
	Messages []string
}
