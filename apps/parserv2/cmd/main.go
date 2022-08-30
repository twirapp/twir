package main

import (
	"fmt"
	"runtime"
	"tsuwari/parser/internal/commands"
	"tsuwari/parser/internal/config/cfg"
	mynats "tsuwari/parser/internal/config/nats"
	"tsuwari/parser/internal/config/redis"
	natshandler "tsuwari/parser/internal/handlers/nats"
	"tsuwari/parser/internal/variables"

	testproto "tsuwari/parser/internal/proto"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/protobuf"
	proto "google.golang.org/protobuf/proto"
)

func main() {
	cfg, err := cfg.New()
	if err != nil || cfg == nil {
		panic("Cannot load config of application")
	}

	r := redis.New(cfg.RedisUrl)
	n, err := mynats.New(cfg.NatsUrl)
	natsJson, err := nats.NewEncodedConn(n, protobuf.PROTOBUF_ENCODER)

	if err != nil {
		panic(err)
	}

	variablesService := variables.New(r)
	commandsService := commands.New(r, variablesService)
	natsHandler := natshandler.New(r, variablesService, commandsService)

	if err != nil {
		panic(err)
	}

	/* natsJson.Subscribe("proto", func(m *nats.Msg) {
		fmt.Println(m.Data)
		m.Respond([]byte(m.Reply))
	}) */

	natsJson.Subscribe("parser.handleProcessCommand", func(m *nats.Msg) {
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
	})

	fmt.Println("Started")

	runtime.Goexit()
	defer r.Close()
	defer n.Close()
}
