package clients

import (
	"fmt"
	"log"

	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/discord"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

func NewDiscord(env string) discord.DiscordClient {
	serverAddress := createClientAddr(env, "discord", constants.DISCORD_SERVER_PORT)

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
	c := discord.NewDiscordClient(conn)
	return c
}
