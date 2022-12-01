package cfg

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RedisUrl                 string  `required:"true"  default:"redis://localhost:6379/0"    envconfig:"REDIS_URL"`
	TwitchClientId           string  `required:"true"                                        envconfig:"TWITCH_CLIENTID"`
	TwitchClientSecret       string  `required:"true"                                        envconfig:"TWITCH_CLIENTSECRET"`
	TwitchCallbackUrl        string  `required:"true"  default:"http://localhost:3005/login" envconfig:"TWITCH_CALLBACKURL"`
	DatabaseUrl              string  `required:"true"                                        envconfig:"DATABASE_URL"`
	AppEnv                   string  `required:"true"  default:"development"                 envconfig:"APP_ENV"`
	SentryDsn                string  `required:"false"                                       envconfig:"SENTRY_DSN"`
	FeedbackTelegramBotToken *string `required:"false"                                       envconfig:"FEEDBACK_TELEGRAM_BOT_TOKEN"`
	FeedbackTelegramUserID   *string `required:"false"                                       envconfig:"FEEDBACK_TELEGRAM_USERID"`
	JwtAccessSecret          string  `required:"false" default:"CoolSecretForAccess"         envconfig:"JWT_ACCESS_SECRET"`
	JwtRefreshSecret         string  `required:"false" default:"CoolSecretForRefresh"        envconfig:"JWT_REFRESH_SECRET"`
}

func New() (*Config, error) {
	var newCfg Config

	var err error

	_ = godotenv.Load(".env")

	if err = envconfig.Process("", &newCfg); err != nil {
		return nil, err
	}

	return &newCfg, nil
}
