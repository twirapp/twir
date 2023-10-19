package clients

import (
	"fmt"
	"log"

	"github.com/satont/twir/libs/grpc/generated/eval"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

func NewEval(env string) eval.EvalClient {
	serverAddress := createClientAddr(env, "eval", constants.EVAL_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(
			fmt.Sprintf(
				`{"loadBalancingConfig": [{"%s":{}}]}`,
				roundrobin.Name,
			),
		),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := eval.NewEvalClient(conn)
	return c
}
