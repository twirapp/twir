package cfg

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NatsUrl            string `required:"true" default:"nats://localhost:4222" envconfig:"NATS_URL"`
	RedisUrl           string `required:"true" default:"redis://localhost:6379/0" envconfig:"REDIS_URL"`
	TwitchClientId     string `required:"true" envconfig:"TWITCH_CLIENTID"`
	TwitchClientSecret string `required:"true" envconfig:"TWITCH_CLIENTSECRET"`
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
