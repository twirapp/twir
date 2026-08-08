package donatexintegration

import (
	"context"

	donatexintegrationentity "github.com/twirapp/twir/libs/entities/donatex_integration"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (donatexintegrationentity.Entity, error)
	Update(ctx context.Context, opts UpdateOpts) error
	Delete(ctx context.Context, channelID string) error
	Create(ctx context.Context, opts CreateOpts) error
}

type CreateOpts struct {
	ChannelID     string
	AccessToken   string
	RefreshToken  string
	DonateXUserID string
	Enabled       bool
	UserName      string
	Avatar        string
}

type UpdateOpts struct {
	ChannelID     string
	AccessToken   *string
	RefreshToken  *string
	DonateXUserID *string
	Enabled       *bool
	UserName      *string
	Avatar        *string
}
