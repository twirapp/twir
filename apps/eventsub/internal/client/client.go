package client

import (
	"context"
	eventsub_framework "github.com/dnsge/twitch-eventsub-framework"
	"github.com/satont/tsuwari/apps/eventsub/internal/creds"
	"github.com/satont/tsuwari/apps/eventsub/internal/types"
)

func NewClient(
	ctx context.Context,
	services *types.Services,
) (*eventsub_framework.SubClient, error) {
	eventSubCreds := creds.NewCreds(ctx, services.Config, services.Grpc.Tokens)
	client := eventsub_framework.NewSubClient(eventSubCreds)

	if services.Config.AppEnv != "production" {
		subs, err := client.GetSubscriptions(ctx, eventsub_framework.StatusAny)
		if err != nil {
			return nil, err
		}

		for _, sub := range subs.Data {
			err = client.Unsubscribe(ctx, sub.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	return client, nil
}
