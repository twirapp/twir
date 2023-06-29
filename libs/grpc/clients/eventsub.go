package clients

import (
	"fmt"
	"log"

	"github.com/satont/twir/libs/grpc/generated/eventsub"
	"github.com/satont/twir/libs/grpc/servers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

func NewEventSub(env string) eventsub.EventSubClient {
	serverAddress := createClientAddr(env, "eventsub", servers.EVENTSUB_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]}`, roundrobin.Name)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := eventsub.NewEventSubClient(conn)
	return c
}
