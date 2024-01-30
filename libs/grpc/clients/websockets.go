package clients

import (
	"log"

	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/websockets"

	"google.golang.org/grpc"
)

func NewWebsocket(env string) websockets.WebsocketClient {
	serverAddress := createClientAddr(env, "websockets", constants.WEBSOCKET_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		defaultClientsOptions...,
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := websockets.NewWebsocketClient(conn)

	return c
}
