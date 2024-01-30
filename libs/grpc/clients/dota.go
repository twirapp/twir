package clients

import (
	"log"

	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/dota"
	"google.golang.org/grpc"
)

func NewDota(env string) dota.DotaClient {
	serverAddress := createClientAddr(env, "dota", constants.DOTA_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		defaultClientsOptions...,
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := dota.NewDotaClient(conn)
	return c
}
