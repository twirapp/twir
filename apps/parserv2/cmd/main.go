package main

import (
	"fmt"
	"tsuwari/parser/internal/config/cfg"
	n "tsuwari/parser/internal/config/nats"
	"tsuwari/parser/internal/config/redis"
	"tsuwari/parser/internal/handlers"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables"

	"github.com/nats-io/nats.go"
)

func main() {
	err := cfg.LoadConfig()
	if err != nil {
		panic("Cannot load config of application")
	}

	redis.Connect()
	n.Connect()
	variables.SetVariables()
	
	v := variables.ParseVariables("$(sender) test $(random|1-100)")
	fmt.Println("variables.ParseVariables:", v)

	cmds, _ := handlers.GetChannelCommands("123")
	fmt.Println("handlers.GetChannelCommands:", cmds)
	cmd := handlers.FindCommandByMessage("!f qweqwe", cmds)
	fmt.Println("handlers.FindCommandByMessage:", cmd)

	user := types.UserInfo{
		UserId: "1",
		UserName: nil,
		UserDisplayName: nil,
		Badges: []string{"MODERATOR"},
	}

	res := handlers.UserHasPermissionToCommand(user.Badges, "MODERATOR")
	fmt.Println("handlers.UserHasPermissionToCommand:", res)

	n.Nats.Subscribe("request", func(m *nats.Msg) {
    m.Respond([]byte("answer is 42"))
	})

	defer redis.Rdb.Close()
}