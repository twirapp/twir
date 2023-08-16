package queue

import (
	"context"
	"github.com/satont/twir/apps/timers/internal/repositories/channels"
	"github.com/satont/twir/apps/timers/internal/repositories/streams"
	"github.com/satont/twir/apps/timers/internal/repositories/timers"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
	"time"
)

type Queue struct {
	timers map[string]*timer
	logger logger.Logger
	config cfg.Config

	timersRepository   timers.Repository
	channelsRepository channels.Repository
	streamsRepository  streams.Repository

	botsGrpc   bots.BotsClient
	parserGrpc parser.ParserClient
}

type timer struct {
	timers.Timer
	LastResponse int

	ticker    *time.Ticker
	doneChann chan bool
}

type Opts struct {
	fx.In

	Lc     fx.Lifecycle
	Logger logger.Logger
	Config cfg.Config

	TimersRepository   timers.Repository
	ChannelsRepository channels.Repository
	StreamsRepository  streams.Repository

	BotsGrpc   bots.BotsClient
	ParserGrpc parser.ParserClient
}

func New(opts Opts) *Queue {
	q := &Queue{
		timers:             make(map[string]*timer),
		logger:             opts.Logger,
		config:             opts.Config,
		timersRepository:   opts.TimersRepository,
		channelsRepository: opts.ChannelsRepository,
		streamsRepository:  opts.StreamsRepository,
		botsGrpc:           opts.BotsGrpc,
		parserGrpc:         opts.ParserGrpc,
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return q.init()
			},
			OnStop: nil,
		},
	)

	return q
}

func (c *Queue) init() error {
	data, err := c.timersRepository.GetAll()
	if err != nil {
		return err
	}

	for _, t := range data {
		if err = c.Add(t.ID); err != nil {
			return err
		}
	}

	return nil
}
