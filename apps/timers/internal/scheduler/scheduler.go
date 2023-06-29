package scheduler

import (
	"fmt"
	"time"

	"github.com/satont/twir/apps/timers/internal/handler"
	"github.com/satont/twir/apps/timers/internal/types"

	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/parser"

	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Scheduler struct {
	internalScheduler *gocron.Scheduler
	cfg               *cfg.Config
	db                *gorm.DB
	logger            *zap.Logger
	Timers            types.Store
	handler           *handler.Handler
}

func New(
	cfg *cfg.Config,
	db *gorm.DB,
	logger *zap.Logger,
	parserGrpc parser.ParserClient,
	botsGrpcClient bots.BotsClient,
) *Scheduler {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.StartAsync()
	store := make(types.Store)
	return &Scheduler{
		internalScheduler: scheduler,
		cfg:               cfg,
		db:                db,
		logger:            logger,
		Timers:            store,
		handler:           handler.New(db, logger, store, parserGrpc, botsGrpcClient),
	}
}

func (c *Scheduler) AddTimer(timer *types.Timer) error {
	c.internalScheduler.RemoveByTag(timer.Model.ID)

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
		c.logger.Sugar().Error(err)
		return err
	}

	c.logger.Info(
		fmt.Sprintf(
			"Queued timer %s#%s for %s channel.",
			timer.Model.Name,
			timer.Model.ID,
			timer.Model.ChannelID,
		),
	)

	return nil
}

func (c *Scheduler) RemoveTimer(id string) error {
	err := c.internalScheduler.RemoveByTag(id)

	c.logger.Sugar().Info(
		fmt.Sprintf(
			"Removed timer %s.",
			id,
		),
	)

	delete(c.Timers, id)

	if err != nil {
		return err
	}

	return nil
}
