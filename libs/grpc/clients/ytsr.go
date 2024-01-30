package clients

import (
	"log"

	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/ytsr"

	"google.golang.org/grpc"
)

func NewYtsr(env string) ytsr.YtsrClient {
	serverAddress := createClientAddr(env, "ytsr", constants.YTSR_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		defaultClientsOptions...,
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := ytsr.NewYtsrClient(conn)

	return c
}
