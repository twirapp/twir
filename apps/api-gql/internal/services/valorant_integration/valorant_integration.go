package valorantintegration

import (
	"context"
	"fmt"

	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/integrations/valorant"
	channelsintegrationsvalorant "github.com/twirapp/twir/libs/repositories/channels_integrations_valorant"
	channelsintegrationsvalorantmodel "github.com/twirapp/twir/libs/repositories/channels_integrations_valorant/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	ChannelsIntegrationsValorantRepository channelsintegrationsvalorant.Repository
	Config                                 cfg.Config
	HenrikApi                              *valorant.HenrikValorantApiClient
}

func New(opts Opts) *Service {
	return &Service{
		repo:      opts.ChannelsIntegrationsValorantRepository,
		henrikApi: opts.HenrikApi,
	}
}

type Service struct {
	repo      channelsintegrationsvalorant.Repository
	henrikApi *valorant.HenrikValorantApiClient
}

func (c *Service) GetChannelStoredMatchesByChannelID(
	ctx context.Context,
	channelID string,
) (*valorant.StoredMatchesResponse, error) {
	integration, err := c.repo.GetByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}
	if integration == channelsintegrationsvalorantmodel.Nil {
		return nil, fmt.Errorf("no valorant integration found for channel id %s", channelID)
	}

	if integration.Data == nil || integration.Data.ValorantPuuid == nil || integration.Data.ValorantActiveRegion == nil {
		return nil, fmt.Errorf("valorant integration data is incomplete for channel id %s", channelID)
	}

	response, err := c.henrikApi.GetProfileStoredMatches(
		ctx,
		*integration.Data.ValorantActiveRegion,
		*integration.Data.ValorantPuuid,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Service) GetChannelMmr(ctx context.Context, channelID string) (
	*valorant.MmrResponse,
	error,
) {
	integration, err := c.repo.GetByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}
	if integration == channelsintegrationsvalorantmodel.Nil {
		return nil, fmt.Errorf("no valorant integration found for channel id %s", channelID)
	}

	if integration.Data == nil || integration.Data.ValorantPuuid == nil || integration.Data.ValorantActiveRegion == nil {
		return nil, fmt.Errorf("valorant integration data is incomplete for channel id %s", channelID)
	}

	response, err := c.henrikApi.GetValorantProfileMmr(
		ctx,
		"pc",
		*integration.Data.ValorantActiveRegion,
		*integration.Data.ValorantPuuid,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}
