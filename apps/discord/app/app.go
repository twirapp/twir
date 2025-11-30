package app

import (
	"github.com/twirapp/twir/apps/discord/internal/discord_go"
	"github.com/twirapp/twir/apps/discord/internal/grpc"
	"github.com/twirapp/twir/apps/discord/internal/messages_updater"
	"github.com/twirapp/twir/apps/discord/internal/sended_messages_store"
	"github.com/twirapp/twir/libs/baseapp"
	channelsintegrationsdiscord "github.com/twirapp/twir/libs/repositories/channels_integrations_discord"
	channelsintegrationsdiscordpostgres "github.com/twirapp/twir/libs/repositories/channels_integrations_discord/datasource/postgres"
	discordsendednotifications "github.com/twirapp/twir/libs/repositories/discord_sended_notifications"
	discordsendednotificationspgx "github.com/twirapp/twir/libs/repositories/discord_sended_notifications/pgx"
	"go.uber.org/fx"
)

var App = fx.Module(
	"discord",
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "discord"}),
	fx.Provide(
		fx.Annotate(
			channelsintegrationsdiscordpostgres.NewFx,
			fx.As(new(channelsintegrationsdiscord.Repository)),
		),
		fx.Annotate(
			discordsendednotificationspgx.NewFx,
			fx.As(new(discordsendednotifications.Repository)),
		),
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
