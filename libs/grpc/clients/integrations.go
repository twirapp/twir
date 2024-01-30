package clients

import (
	"log"

	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/integrations"
	"google.golang.org/grpc"
)

func NewIntegrations(env string) integrations.IntegrationsClient {
	serverAddress := createClientAddr(env, "integrations", constants.INTEGRATIONS_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		defaultClientsOptions...,
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := integrations.NewIntegrationsClient(conn)
	return c
}
