package grpc_clients

import (
	config "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/clients"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"
	"github.com/satont/tsuwari/libs/grpc/generated/eventsub"
	"github.com/satont/tsuwari/libs/grpc/generated/integrations"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
	"github.com/satont/tsuwari/libs/grpc/generated/scheduler"
	"github.com/satont/tsuwari/libs/grpc/generated/timers"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
)

type GrpcClients struct {
	Integrations integrations.IntegrationsClient
	Parser       parser.ParserClient
	EventSub     eventsub.EventSubClient
	Scheduler    scheduler.SchedulerClient
	Timers       timers.TimersClient
	Bots         bots.BotsClient
	Tokens       tokens.TokensClient
}

func NewGrpcClients(cfg *config.Config) *GrpcClients {
	return &GrpcClients{
		Integrations: clients.NewIntegrations(cfg.AppEnv),
		Parser:       clients.NewParser(cfg.AppEnv),
		EventSub:     clients.NewEventSub(cfg.AppEnv),
		Scheduler:    clients.NewScheduler(cfg.AppEnv),
		Timers:       clients.NewTimers(cfg.AppEnv),
		Bots:         clients.NewBots(cfg.AppEnv),
		Tokens:       clients.NewTokens(cfg.AppEnv),
	}
}
