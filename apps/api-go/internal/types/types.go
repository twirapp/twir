package types

import (
	"tsuwari/twitch"

	cfg "tsuwari/config"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/storage/redis"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Services struct {
	DB                  *gorm.DB
	RedisStorage        *redis.Storage
	Validator           *validator.Validate
	ValidatorTranslator ut.Translator
	Twitch              *twitch.Twitch
	Logger              *zap.Logger
	Cfg                 *cfg.Config
	Nats                *nats.Conn
	TgBotApi            *tgbotapi.BotAPI
}
