package clients

import (
	"fmt"
	"github.com/satont/twir/libs/grpc/generated/ytsr"
	"log"

	"github.com/satont/twir/libs/grpc/servers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

func NewYtsr(env string) ytsr.YtsrClient {
	serverAddress := createClientAddr(env, "ytsr", servers.YTSR_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]}`, roundrobin.Name)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := ytsr.NewYtsrClient(conn)

	return c
}
