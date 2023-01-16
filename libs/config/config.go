package cfg

import (
	"os"
	"path/filepath"
	"strings"

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
	HostName                 string  `required:"false" default:"localhost:3005" envconfig:"HOSTNAME"`
	TokensCipherKey          string  `required:"false" default:"pnyfwfiulmnqlhkvixaeligpprcnlyke" envconfig:"TOKENS_CIPHER_KEY"`
}

func New() (*Config, error) {
	var newCfg Config

	var err error

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(wd, "/workspace") {
		wd = "/workspace"
	} else {
		wd = filepath.Join(wd, "..", "..")
	}

	envPath := filepath.Join(wd, ".env")
	_ = godotenv.Load(envPath)

	if err = envconfig.Process("", &newCfg); err != nil {
		return nil, err
	}

	return &newCfg, nil
}
