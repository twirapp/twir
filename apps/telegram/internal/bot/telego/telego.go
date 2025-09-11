package telego

import (
	"context"
	"fmt"

	"github.com/kr/pretty"
	"github.com/mymmrac/telego"
	"github.com/twirapp/twir/apps/telegram/internal/bot"
	"github.com/twirapp/twir/apps/telegram/internal/bot/telego/logger"
	config "github.com/twirapp/twir/libs/config"
	twirlogger "github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Cfg    config.Config
	Logger twirlogger.Logger
}

func New(opts Opts) (*Service, error) {
	if opts.Cfg.TelegramBotToken == "" {
		return nil, fmt.Errorf("telegram bot token is not provided")
	}

	teleBot, err := telego.NewBot(
		opts.Cfg.TelegramBotToken,
		telego.WithLogger(logger.New(opts.Logger)),
	)
	if err != nil {
		return nil, err
	}

	s := &Service{
		b: teleBot,
	}

	botCtx, closeBotCtx := context.WithCancel(context.Background())

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				me, err := teleBot.GetMe(ctx)
				if err != nil {
					return err
				}

				opts.Logger.Info("Telegram bot started as https://t.me/" + me.Username)

				go func() {
					if err := s.Start(botCtx); err != nil {
						panic(err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				closeBotCtx()
				return s.Stop(ctx)
			},
		},
	)

	return s, nil
}

var _ bot.Bot = (*Service)(nil)

type Service struct {
	b *telego.Bot
}

func (s *Service) Start(ctx context.Context) error {
	updates, _ := s.b.UpdatesViaLongPolling(ctx, nil)

	for update := range updates {
		pretty.Println(update.Message.Chat.ID)
	}

	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	return s.b.Close(ctx)
}
