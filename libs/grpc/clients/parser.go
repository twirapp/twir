package clients

import (
	"log"

	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/parser"
	"google.golang.org/grpc"
)

func NewParser(env string) parser.ParserClient {
	serverAddress := createClientAddr(env, "parser", constants.PARSER_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		defaultClientsOptions...,
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := parser.NewParserClient(conn)
	return c
}
