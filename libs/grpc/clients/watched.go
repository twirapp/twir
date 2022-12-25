package clients

import (
	"fmt"
	"log"

	"github.com/satont/tsuwari/libs/grpc/generated/watched"
	"github.com/satont/tsuwari/libs/grpc/servers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

func NewWatched(env string) watched.WatchedClient {
	serverAddress := createClientAddr(env, "watched", servers.WATCHED_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]}`, roundrobin.Name)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := watched.NewWatchedClient(conn)

	return c
}
