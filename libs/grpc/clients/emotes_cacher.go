package clients

import (
	"fmt"
	"log"

	"github.com/satont/twir/libs/grpc/generated/emotes_cacher"

	"github.com/satont/twir/libs/grpc/servers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

func NewEmotesCacher(env string) emotes_cacher.EmotesCacherClient {
	serverAddress := createClientAddr(env, "emotes-cacher", servers.EMOTES_CACHER_SERVER_PORT)

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
	c := emotes_cacher.NewEmotesCacherClient(conn)
	return c
}
