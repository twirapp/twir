package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/satont/tsuwari/apps/timers/internal/scheduler"
	"github.com/satont/tsuwari/apps/timers/internal/types"

	model "github.com/satont/tsuwari/libs/gomodels"

	cfg "github.com/satont/tsuwari/libs/config"

	twitch "github.com/satont/tsuwari/libs/twitch"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/protobuf"
	"github.com/satont/go-helix/v2"
	natstimers "github.com/satont/tsuwari/libs/nats/timers"
)

func main() {
	cfg, err := cfg.New()
	if err != nil || cfg == nil {
		fmt.Println(err)
		panic("Cannot load config of application")
	}

	logger, _ := zap.NewDevelopment()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Sugar().Error(err)
		panic("failed to connect database")
	}

	n, err := nats.Connect(cfg.NatsUrl,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(time.Second),
		nats.Name("Parser-go"),
		nats.Timeout(5*time.Second),
	)
	natsProtoConn, err := nats.NewEncodedConn(n, protobuf.PROTOBUF_ENCODER)

	t := twitch.NewClient(&helix.Options{
		ClientID:     cfg.TwitchClientId,
		ClientSecret: cfg.TwitchClientSecret,
	})

	scheduler := scheduler.New(cfg, t, n, db, logger)

	timers := []*model.ChannelsTimers{}
	err = db.Model(&model.ChannelsTimers{}).
		Where("1 = 1").
		Update("lastTriggerMessageNumber", 0).
		Error
	if err != nil {
		logger.Sugar().Error(err)
	}
	err = db.Preload("Responses").Find(&timers).Error

	if err != nil {
		panic(err)
	} else {
		for _, timer := range timers {
			if timer.Enabled {
				AddTimerByModel(scheduler, timer)
			}
		}
	}

	natsProtoConn.Subscribe("addTimerToQueue", func(m *nats.Msg) {
		data := natstimers.AddTimerToQueue{}
		if err := proto.Unmarshal(m.Data, &data); err != nil {
			logger.Sugar().Error(err)
			return
		}
		timer := &model.ChannelsTimers{}
		if err = db.Where(`"id" = ?`, data.TimerId).Preload("Responses").Take(timer).Error; err != nil {
			logger.Sugar().Error(err)
			return
		}

		AddTimerByModel(scheduler, timer)
		bytes, _ := proto.Marshal(&natstimers.Empty{})
		m.Respond(bytes)
	})

	natsProtoConn.Subscribe("removeTimerFromQueue", func(m *nats.Msg) {
		data := natstimers.RemoveTimerFromQueue{}
		if err := proto.Unmarshal(m.Data, &data); err != nil {
			logger.Sugar().Error(err)
			return
		}

		scheduler.RemoveTimer(data.TimerId)

		bytes, _ := proto.Marshal(&natstimers.Empty{})
		m.Respond(bytes)
	})

	logger.Sugar().Info("Started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	natsProtoConn.Close()
	n.Close()
	log.Fatalf("Exiting")
}

func AddTimerByModel(s *scheduler.Scheduler, t *model.ChannelsTimers) {
	s.AddTimer(&types.Timer{
		Model:     t,
		SendIndex: 0,
	})
}
