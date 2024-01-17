package clients

import (
	"fmt"
	"log"

	"github.com/twirapp/twir/libs/grpc/bots"
	"github.com/twirapp/twir/libs/grpc/constants"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

func NewBots(env string) bots.BotsClient {
	serverAddress := createClientAddr(env, "bots", constants.BOTS_SERVER_PORT)

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
	c := bots.NewBotsClient(conn)
	return c
}
