package channels_categories_aliases

import (
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/twirapp/twir/libs/repositories/channels_categories_aliases/model"
)

type Repository interface {
	GetManyByChannelID(ctx context.Context, channelID string) ([]model.ChannelCategoryAliase, error)
	Create(ctx context.Context, input CreateInput) error
	Delete(ctx context.Context, id ulid.ULID) error
}

type CreateInput struct {
	ChannelID  string
	Alias      string
	CategoryID string
}
