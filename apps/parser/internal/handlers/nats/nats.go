package natshandler

import (
	"tsuwari/parser/internal/commands"
	"tsuwari/parser/internal/variables"

	"github.com/go-redis/redis/v9"
)

type NatsServiceImpl struct {
	redis     *redis.Client
	variables variables.Variables
	commands  commands.Commands
}

type NatsService struct {
	Redis *redis.Client
	Variables variables.Variables
	Commands commands.Commands
}

func New(opts NatsService) NatsServiceImpl {
	return NatsServiceImpl{
		redis: opts.Redis,
		variables: opts.Variables,
		commands: opts.Commands,
	}
}
