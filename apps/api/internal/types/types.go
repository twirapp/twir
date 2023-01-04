package types

import (
	"github.com/satont/tsuwari/libs/twitch"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/satont/tsuwari/apps/api/internal/services/redis"
	"gorm.io/gorm"
)

type Services struct {
	DB                  *gorm.DB
	RedisStorage        *redis.RedisStorage
	Validator           *validator.Validate
	ValidatorTranslator ut.Translator
	Twitch              *twitch.Twitch
	TgBotApi            *tgbotapi.BotAPI
}

type JSONResult struct{}

type DOCApiBadRequest struct {
	Messages string
}

type DOCApiValidationError struct {
	Messages []string
}

type DOCApiInternalError struct {
	Messages []string
}
