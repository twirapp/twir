package clients

import (
	"fmt"
	"log"

	"github.com/satont/tsuwari/libs/grpc/generated/dota"
	"github.com/satont/tsuwari/libs/grpc/servers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

func NewDota(env string) dota.DotaClient {
	serverAddress := createClientAddr(env, "dota", servers.DOTA_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]}`, roundrobin.Name)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := dota.NewDotaClient(conn)
	return c
}
