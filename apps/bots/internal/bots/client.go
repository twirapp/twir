package bots

import (
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/bots/internal/chat_client"
	"github.com/satont/twir/libs/grpc/generated/events"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"github.com/satont/twir/libs/logger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/parser"

	model "github.com/satont/twir/libs/gomodels"

	"gorm.io/gorm"
)

type ClientOpts struct {
	DB             *gorm.DB
	Cfg            cfg.Config
	Logger         logger.Logger
	Model          *model.Bots
	ParserGrpc     parser.ParserClient
	TokensGrpc     tokens.TokensClient
	EventsGrpc     events.EventsClient
	WebsocketsGrpc websockets.WebsocketClient
	Redis          *redis.Client
}

func newBot(opts ClientOpts) *chat_client.ChatClient {
	client := chat_client.New(
		chat_client.Opts{
			DB:             opts.DB,
			Cfg:            opts.Cfg,
			Logger:         opts.Logger,
			Model:          opts.Model,
			ParserGrpc:     opts.ParserGrpc,
			TokensGrpc:     opts.TokensGrpc,
			EventsGrpc:     opts.EventsGrpc,
			WebsocketsGrpc: opts.WebsocketsGrpc,
			Redis:          opts.Redis,
		},
	)

	messagesCounter := promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "bots_messages_counter",
			Help: "The total number of processed messages",
			ConstLabels: prometheus.Labels{
				"botName": client.TwitchUser.Login,
				"botId":   client.TwitchUser.ID,
			},
		},
	)

	prometheus.Register(messagesCounter)

	return client
}
