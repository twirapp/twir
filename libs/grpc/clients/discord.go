package clients

import (
	"log"

	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/discord"
	"google.golang.org/grpc"
)

func NewDiscord(env string) discord.DiscordClient {
	serverAddress := createClientAddr(env, "discord", constants.DISCORD_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		defaultClientsOptions...,
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := discord.NewDiscordClient(conn)
	return c
}
