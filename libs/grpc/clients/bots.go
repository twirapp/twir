package clients

import (
	"log"

	"github.com/twirapp/twir/libs/grpc/bots"
	"github.com/twirapp/twir/libs/grpc/constants"
	"google.golang.org/grpc"
)

func NewBots(env string) bots.BotsClient {
	serverAddress := createClientAddr(env, "bots", constants.BOTS_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		defaultClientsOptions...,
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := bots.NewBotsClient(conn)
	return c
}
