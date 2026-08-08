package donatepay_integration

import (
	"context"
	"errors"
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/integrations"
	donatepayintegration "github.com/twirapp/twir/libs/repositories/donatepay_integration"
	"github.com/twirapp/twir/libs/repositories/donatepay_integration/model"
)

func New(repo donatepayintegration.Repository, twirBus *buscore.Bus) *Service {
	return &Service{
		repo:    repo,
		twirBus: twirBus,
	}
}

type Service struct {
	repo    donatepayintegration.Repository
	twirBus *buscore.Bus
}

var ErrNotFound = errors.New("donatepay integration not found")

func (c *Service) mapModelToEntity(m model.DonatePayIntegration) entity.DonatePayIntegration {
	return entity.DonatePayIntegration{
		ID:        m.ID,
		ChannelID: m.ChannelID,
		ApiKey:    m.ApiKey,
		Enabled:   m.Enabled,
	}
}

func (c *Service) GetByChannelID(
	ctx context.Context,
	channelID string,
) (entity.DonatePayIntegration, error) {
	data, err := c.repo.GetByChannelID(ctx, channelID)
	if err != nil {
		if errors.Is(err, donatepayintegration.ErrNotFound) {
			return entity.DonatePayIntegration{}, ErrNotFound
		}

		return entity.DonatePayIntegration{}, fmt.Errorf("cannot get donatepay: %w", err)
	}

	return c.mapModelToEntity(data), nil
}

func (c *Service) CreateOrUpdate(
	ctx context.Context,
	channelID, apiKey string,
	enabled bool,
) error {
	data, err := c.repo.CreateOrUpdate(ctx, channelID, apiKey, enabled)
	if err != nil {
		return err
	}

	if enabled {
		err = c.twirBus.Integrations.Add.Publish(
			ctx,
			integrations.Request{
				ID:      data.ID.String(),
				Service: integrations.DonatePay,
			},
		)
		if err != nil {
			return err
		}
	} else {
		err = c.twirBus.Integrations.Remove.Publish(
			ctx,
			integrations.Request{
				ID:      data.ChannelID,
				Service: integrations.DonatePay,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
