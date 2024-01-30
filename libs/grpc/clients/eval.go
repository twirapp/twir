package clients

import (
	"log"

	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/eval"
	"google.golang.org/grpc"
)

func NewEval(env string) eval.EvalClient {
	serverAddress := createClientAddr(env, "eval", constants.EVAL_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		defaultClientsOptions...,
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := eval.NewEvalClient(conn)
	return c
}
