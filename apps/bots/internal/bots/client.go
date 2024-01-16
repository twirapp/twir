package bots

import (
	ratelimiter "github.com/aidenwallis/go-ratelimiting/local"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/bots/internal/chat_client"
	"github.com/satont/twir/apps/bots/pkg/tlds"
	"github.com/satont/twir/libs/grpc/events"
	"github.com/satont/twir/libs/grpc/tokens"
	"github.com/satont/twir/libs/grpc/websockets"
	"github.com/satont/twir/libs/logger"

	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/parser"

	model "github.com/satont/twir/libs/gomodels"

	"gorm.io/gorm"
)

type ClientOpts struct {
	DB              *gorm.DB
	Cfg             cfg.Config
	Logger          logger.Logger
	Model           *model.Bots
	ParserGrpc      parser.ParserClient
	TokensGrpc      tokens.TokensClient
	EventsGrpc      events.EventsClient
	WebsocketsGrpc  websockets.WebsocketClient
	Redis           *redis.Client
	JoinRateLimiter ratelimiter.SlidingWindow
	Tlds            *tlds.TLDS
}

func newBot(opts ClientOpts) *chat_client.ChatClient {
	client := chat_client.New(
		chat_client.Opts{
			DB:              opts.DB,
			Cfg:             opts.Cfg,
			Logger:          opts.Logger,
			Model:           opts.Model,
			ParserGrpc:      opts.ParserGrpc,
			TokensGrpc:      opts.TokensGrpc,
			EventsGrpc:      opts.EventsGrpc,
			WebsocketsGrpc:  opts.WebsocketsGrpc,
			Redis:           opts.Redis,
			JoinRateLimiter: opts.JoinRateLimiter,
			Tlds:            opts.Tlds,
		},
	)

	return client
}
