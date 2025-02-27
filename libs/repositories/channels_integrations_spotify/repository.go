package channels_integrations_spotify

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.ChannelIntegrationSpotify, error)
	Create(ctx context.Context, input CreateInput) (model.ChannelIntegrationSpotify, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type UpdateInput struct {
	AccessToken  *string
	RefreshToken *string
	AvatarURI    *string
	Username     *string
	Scopes       *[]string
}

type CreateInput struct {
	ChannelID    string
	AccessToken  string
	RefreshToken string
	AvatarURI    string
	Username     string
	Scopes       []string
}
