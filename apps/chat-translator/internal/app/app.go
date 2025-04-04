package app

import (
	"github.com/twirapp/twir/apps/chat-translator/internal/messaging/twirbus"
	"github.com/twirapp/twir/apps/chat-translator/internal/services/handle_message"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/cache/chat_translations_settings"
	"go.uber.org/fx"

	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"

	channelschattrenslationsrepository "github.com/twirapp/twir/libs/repositories/chat_translation"
	channelschattrenslationsrepositorypostgres "github.com/twirapp/twir/libs/repositories/chat_translation/datasource/postgres"
)

var App = fx.Module(
	"chat-translator",
	fx.NopLogger,
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "chat-translator"}),
	fx.Provide(
		fx.Annotate(
			channelsrepositorypgx.NewFx,
			fx.As(new(channelsrepository.Repository)),
		),
		fx.Annotate(
			channelschattrenslationsrepositorypostgres.NewFx,
			fx.As(new(channelschattrenslationsrepository.Repository)),
		),
	),
	fx.Provide(
		chat_translations_settings.New,
		handle_message.New,
	),
	fx.Invoke(
		twirbus.New,
	),
)
