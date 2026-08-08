package channels

import (
	"context"

	"github.com/google/uuid"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	"github.com/twirapp/twir/libs/entities/platform"
)

type Repository interface {
	GetMany(ctx context.Context, input GetManyInput) ([]channelentity.Channel, error)
	// GetAllByBindingPlatform returns every channel with a binding for p.
	GetAllByBindingPlatform(ctx context.Context, p platform.Platform) ([]channelentity.Channel, error)
	GetByID(ctx context.Context, channelID uuid.UUID) (channelentity.Channel, error)
	GetByIDs(ctx context.Context, channelIDs []uuid.UUID) ([]channelentity.Channel, error)
	GetByApiKey(ctx context.Context, apiKey string) (channelentity.Channel, error)
	// GetByBindingUserID resolves a channel from a platform-scoped linked user ID.
	GetByBindingUserID(ctx context.Context, p platform.Platform, userID uuid.UUID) (channelentity.Channel, error)
	// GetByPlatformChannelID resolves a channel from a platform-scoped provider channel ID.
	GetByPlatformChannelID(ctx context.Context, p platform.Platform, platformChannelID string) (channelentity.Channel, error)
	GetBySlug(ctx context.Context, opts GetBySlugInput) (channelentity.Channel, error)
	Update(ctx context.Context, channelID uuid.UUID, input UpdateInput) (channelentity.Channel, error)
	Create(ctx context.Context) (channelentity.Channel, error)
}

type UpdateInput struct {
	IsEnabled *bool
	ApiKey    *string
}

type GetManyInput struct {
	Enabled *bool
	PerPage int
	Page    int
}

type GetBySlugInput struct {
	Slug     string
	Platform *platform.Platform
}
