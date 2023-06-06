package clients

import (
	"fmt"
	"github.com/satont/tsuwari/libs/grpc/generated/giveaways"
	"github.com/satont/tsuwari/libs/grpc/servers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func NewGiveaways(env string) giveaways.GiveawaysClient {
	serverAddress := createClientAddr(env, "giveaways", servers.GIVEAWAYS_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]}`, roundrobin.Name)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := giveaways.NewGiveawaysClient(conn)
	return c
}
