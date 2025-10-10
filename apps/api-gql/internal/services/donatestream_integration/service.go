package donatestreamintegration

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	model "github.com/twirapp/twir/libs/gomodels"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Gorm *gorm.DB
	KV   kv.KV
}

func New(opts Opts) *Service {
	return &Service{
		gorm: opts.Gorm,
		kv:   opts.KV,
	}
}

type Service struct {
	gorm *gorm.DB
	kv   kv.KV
}

func (c *Service) GetIDByChannelID(ctx context.Context, channelID string) (*uuid.UUID, error) {
	integration := &model.Integrations{}
	if err := c.gorm.WithContext(ctx).Where(
		"service = ?",
		"DONATE_STREAM",
	).First(integration).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("DONATESTREAM integration not enabled on our side")
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

const donateStreamConfirmationKey = "donate_stream_confirmation"

func (c *Service) PostCode(ctx context.Context, channelID string, secret string) error {
	integration := &model.Integrations{}
	if err := c.gorm.WithContext(ctx).Where(
		"service = ?",
		"DONATE_STREAM",
	).First(integration).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("DONATESTREAM integration not enabled on our side")
		}
		return err
	}

	channelIntegration := &model.ChannelsIntegrations{}
	if err := c.gorm.WithContext(ctx).
		Where(
			`"integrationId" = ? AND "channelId" = ?`,
			integration.ID,
			channelID,
		).
		Preload("Integration").
		First(channelIntegration).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("DONATESTREAM integration not enabled on our side")
		}
		return err
	}

	if err := c.kv.
		Set(
			ctx,
			donateStreamConfirmationKey+channelIntegration.ID,
			secret,
			kvoptions.WithExpire(1*time.Hour),
		); err != nil {
		return err
	}

	return nil
}
