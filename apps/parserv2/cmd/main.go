package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	commands "tsuwari/parser/internal/commands"
	"tsuwari/parser/internal/config/cfg"
	mynats "tsuwari/parser/internal/config/nats"
	"tsuwari/parser/internal/config/redis"
	natshandler "tsuwari/parser/internal/handlers/nats"
	variables "tsuwari/parser/internal/variables"

	"github.com/nats-io/nats.go"
)

func main() {
	cfg, err := cfg.New()
	if err != nil || cfg == nil {
		panic("Cannot load config of application")
	}

	r := redis.New(cfg.RedisUrl)
	n, err := mynats.New(cfg.NatsUrl)
	natsJson, err := nats.NewEncodedConn(n, nats.JSON_ENCODER)

	if err != nil {
		panic(err)
	}

	variablesService := variables.New(r)
	commandsService := commands.New(r, variablesService)
	natsHandler := natshandler.New(r, variablesService, commandsService)

	if err != nil {
		panic(err)
	}

	natsJson.Subscribe("parser.handleProcessCommand", func(m *nats.Msg) {
		r := natsHandler.HandleProcessCommand(m)

		if r != nil {
			/* buf := &bytes.Buffer{}
			err = gob.NewEncoder(buf).Encode(&r) */
			res, _ := json.Marshal(r)

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
