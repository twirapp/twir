package handlers

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/bots/internal/di"
	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/satont/tsuwari/apps/bots/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BotInstance struct {
	BotClient *types.BotClient
	Db        *model.Bots
}

type HandlersOpts struct {
	DB         *gorm.DB
	Logger     *zap.Logger
	Cfg        *cfg.Config
	BotClient  *types.BotClient
	ParserGrpc parser.ParserClient
}

type Handlers struct {
	db         *gorm.DB
	logger     *zap.Logger
	BotClient  *types.BotClient
	cfg        *cfg.Config
	parserGrpc parser.ParserClient
	redis      redis.Client

	greetingsCounter prometheus.Counter
	keywordsCounter  prometheus.Counter
}

func CreateHandlers(opts *HandlersOpts) *Handlers {
	redisClient := do.MustInvoke[redis.Client](di.Provider)

	labels := prometheus.Labels{
		"botId":   opts.BotClient.TwitchUser.ID,
		"botName": opts.BotClient.TwitchUser.Login,
	}
	greetingsCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name:        "bots_greetings_counter",
		Help:        "The total number of processed greetings",
		ConstLabels: labels,
	})
	keywordsCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name:        "bots_keywords_counter",
		Help:        "The total number of processed keywords",
		ConstLabels: labels,
	})

	handlersService := &Handlers{
		db:               opts.DB,
		logger:           opts.Logger,
		BotClient:        opts.BotClient,
		cfg:              opts.Cfg,
		parserGrpc:       opts.ParserGrpc,
		greetingsCounter: greetingsCounter,
		keywordsCounter:  keywordsCounter,
		redis:            redisClient,
	}

	prometheus.Register(greetingsCounter)
	prometheus.Register(keywordsCounter)

	return handlersService
}
