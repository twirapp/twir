package integrations

import (
	"context"
	"fmt"

	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type DataFetcherOpts struct {
	fx.In

	Gorm   *gorm.DB
	Config config.Config
}

func NewIntegrationsDataFetcher(opts DataFetcherOpts) *DataFetcher {
	return &DataFetcher{
		gorm:   opts.Gorm,
		config: opts.Config,
	}
}

type DataFetcher struct {
	gorm   *gorm.DB
	config config.Config
}

func (c *DataFetcher) GetIntegrationData(
	ctx context.Context,
	channelId string,
	service gqlmodel.IntegrationService,
) (gqlmodel.IntegrationData, error) {
	integration, err := c.getChannelIntegration(ctx, channelId, service)
	if err != nil {
		return nil, err
	}

	switch service {
	case gqlmodel.IntegrationServiceLastfm:
		if integration.Data == nil || integration.Data.UserName == nil || integration.Data.
			Avatar == nil {
			return nil, nil
		}

		return &gqlmodel.IntegrationDataLastfm{
			Username: *integration.Data.UserName,
			Avatar:   *integration.Data.Avatar,
		}, nil
	case gqlmodel.IntegrationServiceSpotify:
		if integration.Data == nil {
			return nil, nil
		}

		return &gqlmodel.IntegrationDataSpotify{
			Username: *integration.Data.UserName,
			Avatar:   *integration.Data.Avatar,
		}, nil
	case gqlmodel.IntegrationServiceDonationalerts:
		if integration.Data == nil || integration.Data.Name == nil || integration.Data.Avatar == nil {
			return nil, nil
		}

		return &gqlmodel.IntegrationDataDonationAlerts{
			Username: *integration.Data.Name,
			Avatar:   *integration.Data.Avatar,
		}, nil
	case gqlmodel.IntegrationServiceValorant:
		if integration.Data == nil || integration.Data.UserName == nil {
			return nil, nil
		}

		return &gqlmodel.IntegrationDataValorant{
			Username: *integration.Data.UserName,
		}, nil
	case gqlmodel.IntegrationServiceStreamlabs:
		if integration.Data == nil || integration.Data.UserName == nil || integration.Data.Avatar == nil {
			return nil, nil
		}

		return &gqlmodel.IntegrationDataStreamLabs{
			Username: *integration.Data.UserName,
			Avatar:   *integration.Data.Avatar,
		}, nil
	case gqlmodel.IntegrationServiceVk:
		if integration.Data == nil || integration.Data.UserName == nil || integration.Data.Avatar == nil {
			return nil, nil
		}

		return &gqlmodel.IntegrationDataVk{
			Username: *integration.Data.UserName,
			Avatar:   *integration.Data.Avatar,
		}, nil
	}

	return nil, nil
}

func (c *DataFetcher) getChannelIntegration(
	ctx context.Context,
	dashboardId string,
	service gqlmodel.IntegrationService,
) (*model.ChannelsIntegrations, error) {
	entity := model.ChannelsIntegrations{}
	if err := c.gorm.
		WithContext(ctx).
		Where(`"channelId" = ? AND "Integration"."service" = ?`, dashboardId, service).
		Joins("Integration").
		First(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("integration %s not found: %w", service, err)
	}

	return &entity, nil
}
