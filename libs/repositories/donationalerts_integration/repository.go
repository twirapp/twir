package donationalerts_integration

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/donationalerts_integration/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.DonationAlertsIntegration, error)
	Update(ctx context.Context, opts UpdateOpts) error
	Delete(ctx context.Context, channelID string) error
	Create(ctx context.Context, opts CreateOpts) error
}

type CreateOpts struct {
	ChannelID    string
	AccessToken  string
	RefreshToken string
	Enabled      bool
	UserName     string
	Avatar       string
}

type UpdateOpts struct {
	ChannelID    string
	AccessToken  *string
	RefreshToken *string
	Enabled      *bool
	UserName     *string
	Avatar       *string
}
