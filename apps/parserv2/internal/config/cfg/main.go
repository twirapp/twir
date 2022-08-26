package cfg

import (
	"github.com/spf13/viper"
)

type Config struct {
	NatsUrl *string `mapstructure:"NATS_URL"`
	RedisUrl *string `mapstructure:"REDIS_URL"`
}

var Cfg *Config

func LoadConfig(path string) (config Config, err error) {
	
	viper.AddConfigPath("../../")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	viper.AutomaticEnv()
	if err != nil {
			return
	}

	err = viper.Unmarshal(&config)

	Cfg = &config

	return
}