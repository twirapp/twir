package handlers

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"
	cfg "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/gopool"
	"github.com/satont/twir/libs/grpc/generated/events"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/logger"

	"github.com/satont/twir/apps/bots/types"
	"gorm.io/gorm"
)

type BotInstance struct {
	BotClient *types.BotClient
	Db        *model.Bots
}

type Opts struct {
	DB         *gorm.DB
	Logger     logger.Logger
	Cfg        cfg.Config
	BotClient  *types.BotClient
	ParserGrpc parser.ParserClient
	EventsGrpc events.EventsClient
	TokensGrpc tokens.TokensClient
	Redis      *redis.Client
}

type Handlers struct {
	db         *gorm.DB
	logger     logger.Logger
	BotClient  *types.BotClient
	cfg        cfg.Config
	parserGrpc parser.ParserClient
	eventsGrpc events.EventsClient
	redis      *redis.Client
	tokensGrpc tokens.TokensClient

	workersPool *gopool.Pool

	greetingsCounter prometheus.Counter
	keywordsCounter  prometheus.Counter
}

func CreateHandlers(opts *Opts) *Handlers {
	labels := prometheus.Labels{
		"botId":   opts.BotClient.TwitchUser.ID,
		"botName": opts.BotClient.TwitchUser.Login,
	}
	greetingsCounter := promauto.NewCounter(
		prometheus.CounterOpts{
			Name:        "bots_greetings_counter",
			Help:        "The total number of processed greetings",
			ConstLabels: labels,
		},
	)
	keywordsCounter := promauto.NewCounter(
		prometheus.CounterOpts{
			Name:        "bots_keywords_counter",
			Help:        "The total number of processed keywords",
			ConstLabels: labels,
		},
	)

	workersPool := gopool.NewPool(1000)

	handlersService := &Handlers{
		db:               opts.DB,
		logger:           opts.Logger,
		BotClient:        opts.BotClient,
		cfg:              opts.Cfg,
		parserGrpc:       opts.ParserGrpc,
		greetingsCounter: greetingsCounter,
		keywordsCounter:  keywordsCounter,
		redis:            opts.Redis,
		eventsGrpc:       opts.EventsGrpc,
		workersPool:      workersPool,
		tokensGrpc:       opts.TokensGrpc,
	}

	prometheus.Register(greetingsCounter)
	prometheus.Register(keywordsCounter)

	return handlersService
}
