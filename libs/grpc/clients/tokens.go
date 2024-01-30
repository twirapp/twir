package clients

import (
	"log"

	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/tokens"

	"google.golang.org/grpc"
)

func NewTokens(env string) tokens.TokensClient {
	serverAddress := createClientAddr(env, "tokens", constants.TOKENS_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		defaultClientsOptions...,
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := tokens.NewTokensClient(conn)

	return c
}
