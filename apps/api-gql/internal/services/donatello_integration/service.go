package donatellointegration

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/guregu/null"
	model "github.com/twirapp/twir/libs/gomodels"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Gorm *gorm.DB
}

func New(opts Opts) *Service {
	return &Service{
		gorm: opts.Gorm,
	}
}

type Service struct {
	gorm *gorm.DB
}

func (c *Service) GetIDByChannelID(ctx context.Context, channelID string) (*uuid.UUID, error) {
	integration := &model.Integrations{}
	if err := c.gorm.WithContext(ctx).Where(
		"service = ?",
		"DONATELLO",
	).First(integration).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("DONATELLO integration not enabled on our side")
		}
		return nil, err
	}

	channelIntegration := &model.ChannelsIntegrations{}
	if err := c.gorm.WithContext(ctx).
		Where(
			`"integrationId" = ? AND "channelId" = ?`,
			integration.ID,
			channelID,
		).
		Preload("Integration").
		Find(channelIntegration).Error; err != nil {
		return nil, err
	}

	if channelIntegration.ID == "" {
		channelIntegration = &model.ChannelsIntegrations{
			Enabled:       true,
			ChannelID:     channelID,
			IntegrationID: integration.ID,
			AccessToken:   null.String{},
			RefreshToken:  null.String{},
			ClientID:      null.String{},
			ClientSecret:  null.String{},
			APIKey:        null.String{},
			Integration:   integration,
			Data:          &model.ChannelsIntegrationsData{},
		}

		if err := c.gorm.WithContext(ctx).Save(channelIntegration).Error; err != nil {
			return nil, err
		}
	}

	id, err := uuid.Parse(channelIntegration.ID)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
