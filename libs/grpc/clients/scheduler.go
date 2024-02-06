package clients

import (
	"log"

	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/scheduler"
	"google.golang.org/grpc"
)

func NewScheduler(env string) scheduler.SchedulerClient {
	serverAddress := createClientAddr(env, "scheduler", constants.SCHEDULER_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		defaultClientsOptions...,
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := scheduler.NewSchedulerClient(conn)
	return c
}
