package main

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/telegram/internal/bot"
	"github.com/twirapp/twir/apps/telegram/internal/bot/telego"
	"github.com/twirapp/twir/libs/baseapp"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		baseapp.CreateBaseApp(
			baseapp.Opts{
				AppName: "telegram",
			},
		),
		fx.Provide(
			fx.Annotate(
				telego.New,
				fx.As(new(bot.Bot)),
			),
		),
		fx.Invoke(
			func(b bot.Bot) {
				_ = b.SendMessage(
					context.TODO(),
					bot.SendMessageInput{
						ChatID:    210144787,
						Text:      "🚀 Bot started",
						ParseMode: lo.ToPtr(bot.SendMessageParseModeHTML),
					},
				)
			},
		),
	).Run()
}
