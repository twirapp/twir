package integrations

import (
	"context"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_deps"
	model "github.com/satont/tsuwari/libs/gomodels"
)

type Integrations struct {
	*impl_deps.Deps
}

func (c *Integrations) getChannelIntegrationByService(
	ctx context.Context,
	service model.IntegrationService,
) (*model.ChannelsIntegrations, error) {
	integration := &model.Integrations{}
	if err := c.Db.WithContext(ctx).Where("service = ?", service).First(integration).Error; err != nil {
		return nil, err
	}

	channelIntegration := &model.ChannelsIntegrations{}
	if err := c.Db.WithContext(ctx).Where(`"integrationId" = ?`, integration.ID).First(channelIntegration).Error; err != nil {
		return nil, err
	}

	return channelIntegration, nil
}
