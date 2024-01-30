package clients

import (
	"log"

	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/events"

	"google.golang.org/grpc"
)

func NewEvents(env string) events.EventsClient {
	serverAddress := createClientAddr(env, "events", constants.EVENTS_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		defaultClientsOptions...,
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := events.NewEventsClient(conn)
	return c
}
