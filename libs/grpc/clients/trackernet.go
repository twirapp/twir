package clients

import (
	"fmt"
	"github.com/satont/tsuwari/libs/grpc/generated/trackernet"
	"log"

	"github.com/satont/tsuwari/libs/grpc/servers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

func NewTrackernet(env string) trackernet.TrackernetClient {
	serverAddress := createClientAddr(env, "trackernet", servers.TRACKERNET_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]}`, roundrobin.Name)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := trackernet.NewTrackernetClient(conn)

	return c
}
