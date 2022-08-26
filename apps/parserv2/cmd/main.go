package main

import (
	"fmt"
	"tsuwari/parser/internal/config/cfg"
	"tsuwari/parser/internal/config/nats"
	"tsuwari/parser/internal/config/redis"
	"tsuwari/parser/internal/handlers"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables"
)

func main() {
	err := cfg.LoadConfig()
	if err != nil {
		panic("Cannot load config of application")
	}

	redis.Connect()
	nats.Connect()
	variables.SetVariables()
	
	variables.ParseVariables("$(sender) test $(random|1-100)")

	cmds, _ := handlers.GetChannelCommands("123")
	cmd := handlers.FindCommandByMessage("!First qweqwe", cmds)
	fmt.Println("cmd:", cmd)

	userName := "test"
	displayName := "Test"

	user := types.UserInfo{
		UserId: "1",
		UserName: &userName,
		UserDisplayName: &displayName,
		Badges: []string{"MODERATOR"},
	}

	res := handlers.UserHasPermissionToCommand(user.Badges, "MODERATOR")
	fmt.Println(res)

	defer redis.Rdb.Close()
}