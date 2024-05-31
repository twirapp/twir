package app

import (
	"github.com/satont/twir/apps/discord/internal/discord_go"
	"github.com/satont/twir/apps/discord/internal/grpc"
	"github.com/satont/twir/apps/discord/internal/messages_updater"
	"github.com/satont/twir/apps/discord/internal/sended_messages_store"
	cfg "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/libs/baseapp"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
)

var App = fx.Module(
	"discord",
	baseapp.CreateBaseApp("discord"),
	fx.Provide(
		buscore.NewNatsBusFx("discord"),
		sended_messages_store.New,
		messages_updater.New,
		discord_go.New,
		func(config cfg.Config) tokens.TokensClient {
			return clients.NewTokens(config.AppEnv)
		},
		grpc.New,
	),
	fx.Invoke(
		messages_updater.New,
		grpc.New,
	),
)
