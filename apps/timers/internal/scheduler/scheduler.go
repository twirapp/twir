package scheduler

import (
	"context"
	model "github.com/satont/twir/libs/gomodels"
	"go.uber.org/fx"
	"log/slog"
	"time"

	"github.com/satont/twir/apps/timers/internal/handler"
	"github.com/satont/twir/apps/timers/internal/types"

	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/parser"

	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
)

type Scheduler struct {
	internalScheduler *gocron.Scheduler
	cfg               *cfg.Config
	db                *gorm.DB
	Timers            types.Store
	handler           *handler.Handler
}

type Opts struct {
	fx.In

	Lc             fx.Lifecycle
	Cfg            *cfg.Config
	Db             *gorm.DB
	ParserGrpc     parser.ParserClient
	BotsGrpcClient bots.BotsClient
}

func New(opts Opts) *Scheduler {
	scheduler := gocron.NewScheduler(time.UTC)

	store := make(types.Store)
	s := &Scheduler{
		internalScheduler: scheduler,
		cfg:               opts.Cfg,
		db:                opts.Db,
		Timers:            store,
		handler:           handler.New(opts.Db, store, opts.ParserGrpc, opts.BotsGrpcClient, opts.Cfg),
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				slog.Info("invoked")
				return nil
			},
			OnStop: nil,
		},
	)

	var timers []model.ChannelsTimers
	err := opts.Db.Model(&model.ChannelsTimers{}).
		Where("1 = 1").
		Update("lastTriggerMessageNumber", 0).
		Error
	if err != nil {
		slog.Error(err.Error())
	}
	err = opts.Db.Preload("Responses").Preload("Channel").Find(&timers).Error

	if err != nil {
		slog.Error(err.Error())
	} else {
		for _, timer := range timers {
			if timer.Channel != nil && (!timer.Channel.IsEnabled || timer.Channel.IsBanned) || !timer.Enabled {
				continue
			}

			err = s.AddTimer(types.Timer{Model: timer, SendIndex: 0})
			if err != nil {
				slog.Error(err.Error(), "name", timer.Name, "channelId", timer.ChannelID)
			}
		}
	}

	scheduler.StartAsync()

	return s
}

func (c *Scheduler) AddTimer(timer types.Timer) error {
	c.RemoveTimer(timer.Model.ID)

	c.Timers[timer.Model.ID] = &timer

	s := c.internalScheduler.Tag(timer.Model.ID).Every(int(timer.Model.TimeInterval))
	if c.cfg.AppEnv != "production" {
		s = s.Second()
	} else {
		s = s.Minute()
	}

	if _, err := s.DoWithJobDetails(c.handler.Handle); err != nil {
		slog.Error(err.Error())
		return err
	}

	slog.Info("Queued timer", "name", timer.Model.Name, "id", timer.Model.ID, "channelId", timer.Model.ChannelID)

	return nil
}

func (c *Scheduler) RemoveTimer(id string) error {
	err := c.internalScheduler.RemoveByTag(id)

	delete(c.Timers, id)
	if err != nil {
		return err
	}

	slog.Info("Removed timer", "id", id)

	return nil
}
