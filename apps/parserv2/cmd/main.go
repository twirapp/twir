package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	"tsuwari/parser/internal/commands"
	"tsuwari/parser/internal/config/cfg"
	mynats "tsuwari/parser/internal/config/nats"
	"tsuwari/parser/internal/config/redis"
	twitch "tsuwari/parser/internal/config/twitch"
	natshandler "tsuwari/parser/internal/handlers/nats"
	"tsuwari/parser/internal/variables"

	testproto "tsuwari/parser/internal/proto"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/protobuf"
	proto "google.golang.org/protobuf/proto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := cfg.New()
	if err != nil || cfg == nil {
		panic("Cannot load config of application")
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	r := redis.New(cfg.RedisUrl)
	defer r.Close()
	n, err := mynats.New(cfg.NatsUrl)
	natsJson, err := nats.NewEncodedConn(n, protobuf.PROTOBUF_ENCODER)

	if err != nil {
		panic(err)
	}
	defer n.Close()

	twitchClient := twitch.New(*cfg)
	variablesService := variables.New(r, twitchClient, db)
	commandsService := commands.New(r, variablesService, db)
	natsHandler := natshandler.New(r, variablesService, commandsService)

	if err != nil {
		panic(err)
	}

	/* natsJson.Subscribe("proto", func(m *nats.Msg) {
		fmt.Println(m.Data)
		m.Respond([]byte(m.Reply))
	}) */

	natsJson.Subscribe("parser.handleProcessCommand", func(m *nats.Msg) {
		start := time.Now()
		r := natsHandler.HandleProcessCommand(m)

		if r != nil {
			res, _ := proto.Marshal(&testproto.Response{
				Responses: *r,
			})

			if err == nil {
				m.Respond(res)
			} else {
				fmt.Println(err)
			}
		} else {
			m.Respond([]byte{})
		}

		log.Printf("Binomial took %s", time.Since(start))
	})

	fmt.Println("Started")

	// runtime.Goexit()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Fatalf("Exiting")
}
