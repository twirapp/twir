package scheduler

import (
	"context"
	model "github.com/satont/twir/libs/gomodels"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
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

func New(
	cfg *cfg.Config,
	db *gorm.DB,
	parserGrpc parser.ParserClient,
	botsGrpcClient bots.BotsClient,
	lc fx.Lifecycle,
) *Scheduler {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.StartAsync()
	store := make(types.Store)
	s := &Scheduler{
		internalScheduler: scheduler,
		cfg:               cfg,
		db:                db,
		Timers:            store,
		handler:           handler.New(db, store, parserGrpc, botsGrpcClient, cfg),
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				var timers []model.ChannelsTimers
				err := db.Model(&model.ChannelsTimers{}).
					Where("1 = 1").
					Update("lastTriggerMessageNumber", 0).
					Error
				if err != nil {
					return err
				}
				err = db.Preload("Responses").Preload("Channel").Find(&timers).Error

				if err != nil {
					return err
				} else {
					for _, timer := range timers {
						if timer.Channel != nil && (!timer.Channel.IsEnabled || timer.Channel.IsBanned) {
							continue
						}

						if timer.Enabled {
							err = s.AddTimer(&types.Timer{Model: &timer, SendIndex: 0})
							if err != nil {
								slog.Error(err.Error(), "name", timer.Name, "channelId", timer.ChannelID)
							}
						}
					}
				}

				return nil
			},
			OnStop: nil,
		},
	)

	return s
}

func (c *Scheduler) AddTimer(timer *types.Timer) error {
	c.RemoveTimer(timer.Model.ID)

	var unit time.Duration

	if c.cfg.AppEnv != "production" {
		unit = time.Second
	} else {
		unit = time.Minute
	}

	c.Timers[timer.Model.ID] = timer

	_, err := c.internalScheduler.
		Every(int(unit * time.Duration(timer.Model.TimeInterval) / time.Millisecond)).
		Tag(timer.Model.ID).
		Millisecond().
		DoWithJobDetails(c.handler.Handle)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	slog.Info("Queued timer", "name", timer.Model.Name, "id", timer.Model.ID, "channelId", timer.Model.ChannelID)

	return nil
}

func (c *Scheduler) RemoveTimer(id string) error {
	err := c.internalScheduler.RemoveByTag(id)

	if err != nil {
		return err
	}

	slog.Info("Removed timer", "id", id)

	delete(c.Timers, id)

	return nil
}
