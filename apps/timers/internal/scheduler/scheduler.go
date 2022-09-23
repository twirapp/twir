package scheduler

import (
	"fmt"
	"time"
	cfg "tsuwari/config"
	"tsuwari/timers/internal/handler"
	"tsuwari/timers/internal/types"
	"tsuwari/twitch"

	"github.com/go-co-op/gocron"
	redis "github.com/go-redis/redis/v9"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Scheduler struct {
	internalScheduler *gocron.Scheduler
	cfg *cfg.Config
	redis *redis.Client
	handler *handler.Handler
	twitch *twitch.Twitch
	nats *nats.Conn
	db *gorm.DB
	logger *zap.Logger
}

func New(cfg *cfg.Config, redis *redis.Client, twitch *twitch.Twitch, nats *nats.Conn, db *gorm.DB, logger *zap.Logger) *Scheduler {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.StartAsync()

	return &Scheduler{
		internalScheduler: scheduler, 
		cfg: cfg, 
		redis: redis,
		handler: handler.New(redis, twitch, nats, db, logger),
		twitch: twitch,
		nats: nats,
		db: db,
		logger: logger,
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

	_, err := c.internalScheduler.
		Every(int(unit * time.Duration(timer.Model.TimeInterval) / time.Millisecond)).
		Tag(timer.Model.ID).
		Millisecond().
		DoWithJobDetails(c.handler.Handle, timer)

	if err != nil {
		return err
	}

	c.logger.Info(fmt.Sprintf(
		"Queued timer %s#%s for %s channel.",
		timer.Model.Name,
		timer.Model.ID,
		timer.Model.ChannelID,
	))

	return nil
}

func (c *Scheduler) RemoveTimer(id string) error {
	err := c.internalScheduler.RemoveByTag(id)

	fmt.Println(
		fmt.Sprintf(
			"Removed timer %s.",
			id,
		),
	)

	if err != nil {
		return err
	}

	return nil
}