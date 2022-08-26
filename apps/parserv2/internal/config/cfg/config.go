package cfg

import (
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NatsUrl string `required:"true" default:"nats://localhost:4222" envconfig:"NATS_URL"`
	RedisUrl string `required:"true" default:"redis://localhost:6379/0" envconfig:"REDIS_URL"`
}

var (
	once   sync.Once
	Cfg *Config
)

func LoadConfig() error {
	var err error
	once.Do(func() {
		var cfg Config
		// If you run it locally and through terminal please set up this in Load function (../.env)
		_ = godotenv.Load(".env")

		if err = envconfig.Process("", &cfg); err != nil {
			return
		}

		Cfg = &cfg
	})

	return err
}
