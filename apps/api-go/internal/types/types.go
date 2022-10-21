package types

import (
	"tsuwari/twitch"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/storage/redis"
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
}
