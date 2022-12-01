package clients

import (
	"fmt"
	"log"

	"github.com/satont/tsuwari/libs/grpc/generated/integrations"
	"github.com/satont/tsuwari/libs/grpc/servers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

func NewIntegrations(env string) integrations.IntegrationsClient {
	serverAddress := createClientAddr(env, "integrations", servers.INTEGRATIONS_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]}`, roundrobin.Name)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := integrations.NewIntegrationsClient(conn)
	return c
}
