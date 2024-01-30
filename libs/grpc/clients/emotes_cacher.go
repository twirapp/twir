package clients

import (
	"log"

	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/emotes_cacher"

	"google.golang.org/grpc"
)

func NewEmotesCacher(env string) emotes_cacher.EmotesCacherClient {
	serverAddress := createClientAddr(env, "emotes-cacher", constants.EMOTES_CACHER_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		defaultClientsOptions...,
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := emotes_cacher.NewEmotesCacherClient(conn)
	return c
}
