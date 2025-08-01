package donatepayintegration

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/donatepay_integration/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.DonatePayIntegration, error)
	CreateOrUpdate(ctx context.Context, channelID, apiKey string, enabled bool) (model.DonatePayIntegration, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
