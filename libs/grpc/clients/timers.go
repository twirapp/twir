package clients

import (
	"log"

	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/timers"
	"google.golang.org/grpc"
)

func NewTimers(env string) timers.TimersClient {
	serverAddress := createClientAddr(env, "timers", constants.TIMERS_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		defaultClientsOptions...,
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := timers.NewTimersClient(conn)
	return c
}
