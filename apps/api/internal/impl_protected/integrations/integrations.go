package integrations

import (
	"context"

	"github.com/guregu/null"
	"github.com/twirapp/twir/apps/api/internal/impl_deps"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/integrations"
)

type Integrations struct {
	*impl_deps.Deps
}

func (c *Integrations) getIntegrationByService(
	ctx context.Context,
	service model.IntegrationService,
) (*model.Integrations, error) {
	integration := &model.Integrations{}
	if err := c.Db.WithContext(ctx).Where(
		"service = ?",
		service,
	).First(integration).Error; err != nil {
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
			Enabled:       true,
			ChannelID:     channelId,
			IntegrationID: integration.ID,
			AccessToken:   null.String{},
			RefreshToken:  null.String{},
			ClientID:      null.String{},
			ClientSecret:  null.String{},
			APIKey:        null.String{},
			Integration:   integration,
			Data:          &model.ChannelsIntegrationsData{},
		}

		if err := c.Db.WithContext(ctx).Save(channelIntegration).Error; err != nil {
			return nil, err
		}
	}

	return channelIntegration, nil
}

func (c *Integrations) sendBusEvent(ctx context.Context, integrationId string, isAdd bool) error {
	req := integrations.Request{
		ID: integrationId,
	}

	if isAdd {
		return c.Bus.Integrations.Add.Publish(ctx, req)
	} else {
		return c.Bus.Integrations.Remove.Publish(ctx, req)
	}
}
