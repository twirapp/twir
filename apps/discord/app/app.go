package app

import (
	"github.com/satont/twir/apps/discord/internal/discord_go"
	"github.com/satont/twir/apps/discord/internal/grpc"
	"github.com/satont/twir/apps/discord/internal/messages_updater"
	"github.com/satont/twir/apps/discord/internal/sended_messages_store"
	"github.com/twirapp/twir/libs/baseapp"
	"go.uber.org/fx"
)

var App = fx.Module(
	"discord",
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "discord"}),
	fx.Provide(
		sended_messages_store.New,
		messages_updater.New,
		discord_go.New,
		grpc.New,
	),
	fx.Invoke(
		messages_updater.New,
		grpc.New,
	),
)
