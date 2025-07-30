package donatepay_integration

import (
	"context"
	"errors"
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	donatepayintegration "github.com/twirapp/twir/libs/repositories/donatepay_integration"
	"github.com/twirapp/twir/libs/repositories/donatepay_integration/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Repo donatepayintegration.Repository
}

func New(opts Opts) *Service {
	return &Service{
		repo: opts.Repo,
	}
}

type Service struct {
	repo donatepayintegration.Repository
}

var ErrNotFound = errors.New("donatepay integration not found")

func (c *Service) mapModelToEntity(m model.DonatePayIntegration) entity.DonatePayIntegration {
	return entity.DonatePayIntegration{
		ID:        m.ID,
		ChannelID: m.ChannelID,
		ApiKey:    m.ApiKey,
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

func (c *Service) CreateOrUpdate(ctx context.Context, channelID, apiKey string) error {
	_, err := c.repo.CreateOrUpdate(ctx, channelID, apiKey)

	return fmt.Errorf("cannot create donatepay integration: %w", err)
}
