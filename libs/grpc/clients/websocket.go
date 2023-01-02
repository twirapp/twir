package clients

import (
	"fmt"
	"github.com/satont/tsuwari/libs/grpc/generated/websocket"
	"log"

	"github.com/satont/tsuwari/libs/grpc/servers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

func NewWebsocket(env string) websocket.WebsocketClient {
	serverAddress := createClientAddr(env, "websocket", servers.WEBSOCKET_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]}`, roundrobin.Name)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := websocket.NewWebsocketClient(conn)

	return c
}
