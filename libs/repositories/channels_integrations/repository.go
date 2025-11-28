package channels_integrations

import (
"context"

"github.com/twirapp/twir/libs/repositories/channels_integrations/model"
integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
)

type Repository interface {
	GetByChannelAndService(
ctx context.Context,
channelID string,
service integrationsmodel.Service,
) (model.ChannelIntegration, error)
	Create(ctx context.Context, input CreateInput) (model.ChannelIntegration, error)
	Update(ctx context.Context, id string, input UpdateInput) error
}

type CreateInput struct {
	ChannelID     string
	IntegrationID string
	Enabled       bool
	AccessToken   *string
	RefreshToken  *string
	Data          *model.Data
}

type UpdateInput struct {
	Enabled      *bool
	AccessToken  *string
	RefreshToken *string
	Data         *model.Data
}
