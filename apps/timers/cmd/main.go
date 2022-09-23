package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	model "tsuwari/models"
	"tsuwari/timers/internal/scheduler"
	"tsuwari/timers/internal/services/redis"
	"tsuwari/timers/internal/types"

	cfg "tsuwari/config"
	twitch "tsuwari/twitch"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/protobuf"
	natstimers "github.com/satont/tsuwari/nats/timers"
)


func main() {
	cfg, err := cfg.New()
	if err != nil || cfg == nil {
		fmt.Println(err)
		panic("Cannot load config of application")
	}

	var logger *zap.Logger

	if cfg.AppEnv == "development" {
		l, _ := zap.NewDevelopment()
		logger = l
	} else {
		l, _ := zap.NewProduction()
		logger = l
	}
	
	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	r := redis.New(cfg.RedisUrl)

	n, err := nats.Connect(cfg.NatsUrl,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(time.Second),
		nats.Name("Parser-go"),
		nats.Timeout(5*time.Second),
	)
	natsProtoConn, err := nats.NewEncodedConn(n, protobuf.PROTOBUF_ENCODER)

	t := twitch.New(cfg.TwitchClientId, cfg.TwitchClientSecret)
	
	scheduler := scheduler.New(cfg, r, t, n, db, logger)

	timers := []model.ChannelsTimers{}

	err = db.Find(&timers).Error

	if err != nil {
		panic(err)
	} else {
		for _, timer := range timers {
			if timer.Enabled {
				AddTimerByModel(scheduler, &timer)
			}
		}
	}

	natsProtoConn.Subscribe("addTimerToQueue", func(m *nats.Msg) {
		data := natstimers.AddTimerToQueue{}
		if err := proto.Unmarshal(m.Data, &data); err != nil {
			fmt.Println(err)
			return
		}
		timer := &model.ChannelsTimers{}
		if err = db.Where(`"id" = ?`, data.TimerId).Take(timer).Error; err != nil {
			fmt.Println(err)
			return
		}

		AddTimerByModel(scheduler, timer)
		bytes, _ := proto.Marshal(&natstimers.Empty{})
		m.Respond(bytes)
	})

	natsProtoConn.Subscribe("removeTimerFromQueue", func(m *nats.Msg) {
		data := natstimers.RemoveTimerFromQueue{}
		if err := proto.Unmarshal(m.Data, &data); err != nil {
			fmt.Println(err)
			return
		}

		scheduler.RemoveTimer(data.TimerId)

		bytes, _ := proto.Marshal(&natstimers.Empty{})
		m.Respond(bytes)
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	r.Close()
	natsProtoConn.Close()
	n.Close()
	log.Fatalf("Exiting")
}

func AddTimerByModel (s *scheduler.Scheduler, t *model.ChannelsTimers) {
	s.AddTimer(&types.Timer{
		Model: t,
		SendIndex: 0,
	})
}