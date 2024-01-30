package clients

import (
	"log"

	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/eventsub"
	"google.golang.org/grpc"
)

func NewEventSub(env string) eventsub.EventSubClient {
	serverAddress := createClientAddr(env, "eventsub", constants.EVENTSUB_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		defaultClientsOptions...,
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := eventsub.NewEventSubClient(conn)
	return c
}
