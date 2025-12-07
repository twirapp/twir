package faceitintegration

import (
	"context"

	faceitintegrationentity "github.com/twirapp/twir/libs/entities/faceit_integration"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (faceitintegrationentity.Entity, error)
	Update(ctx context.Context, opts UpdateOpts) error
	Delete(ctx context.Context, channelID string) error
	Create(ctx context.Context, opts CreateOpts) error
}

type CreateOpts struct {
	ChannelID   string
	AccessToken string
	Enabled     bool
	UserName    string
	Avatar      string
	Game        string
}

type UpdateOpts struct {
	ChannelID   string
	AccessToken *string
	Enabled     *bool
	UserName    *string
	Avatar      *string
	Game        *string
}
