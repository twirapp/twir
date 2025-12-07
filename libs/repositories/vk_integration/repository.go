package vk_integration

import (
	"context"

	"github.com/twirapp/twir/libs/entities/vk_integration"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (vk_integration.Entity, error)
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
}

type UpdateOpts struct {
	ChannelID   string
	AccessToken *string
	Enabled     *bool
	UserName    *string
	Avatar      *string
}
