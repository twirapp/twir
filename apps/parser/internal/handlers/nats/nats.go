package natshandler

import (
	"tsuwari/parser/internal/commands"
	"tsuwari/parser/internal/variables"

	"github.com/go-redis/redis/v9"
)

type natsService struct {
	redis     *redis.Client
	variables variables.Variables
	commands  commands.Commands
}

func New(redis *redis.Client, variables variables.Variables, commands commands.Commands) natsService {
	return natsService{
		redis,
		variables,
		commands,
	}
}
