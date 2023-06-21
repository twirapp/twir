package integrations

import (
	"context"
	"github.com/guregu/null"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_deps"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/integrations"
)

type Integrations struct {
	*impl_deps.Deps
}

func (c *Integrations) getIntegrationByService(
	ctx context.Context,
	service model.IntegrationService,
) (*model.Integrations, error) {
	integration := &model.Integrations{}
	if err := c.Db.WithContext(ctx).Where("service = ?", service).First(integration).Error; err != nil {
		return nil, err
	}

	return integration, nil
}

func (c *Integrations) getChannelIntegrationByService(
	ctx context.Context,
	service model.IntegrationService,
	channelId string,
) (*model.ChannelsIntegrations, error) {
	integration, err := c.getIntegrationByService(ctx, service)
	if err != nil {
		return nil, err
	}

	channelIntegration := &model.ChannelsIntegrations{}
	if err := c.Db.WithContext(ctx).
		Where(
			`"integrationId" = ? AND "channelId" = ?`,
			integration.ID,
			channelId,
		).
		Preload("Integration").
		Find(channelIntegration).Error; err != nil {
		return nil, err
	}

	if channelIntegration.ID == "" {
		channelIntegration = &model.ChannelsIntegrations{
			Enabled:       false,
			ChannelID:     channelId,
			IntegrationID: integration.ID,
			AccessToken:   null.String{},
			RefreshToken:  null.String{},
			ClientID:      null.String{},
			ClientSecret:  null.String{},
			APIKey:        null.String{},
			Integration:   integration,
		}

		if err := c.Db.WithContext(ctx).Save(channelIntegration).Error; err != nil {
			return nil, err
		}
	}

	return channelIntegration, nil
}

func (c *Integrations) sendGrpcEvent(ctx context.Context, integrationId string, isAdd bool) {
	if isAdd {
		c.Grpc.Integrations.AddIntegration(
			ctx,
			&integrations.Request{
				Id: integrationId,
			},
		)
	} else {
		c.Grpc.Integrations.RemoveIntegration(
			ctx,
			&integrations.Request{
				Id: integrationId,
			},
		)
	}
}
