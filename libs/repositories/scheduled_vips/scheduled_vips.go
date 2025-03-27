package scheduled_vips

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/twirapp/twir/libs/repositories/scheduled_vips/model"
)

type Repository interface {
	GetByID(ctx context.Context, id ulid.ULID) (model.ScheduledVip, error)
	GetMany(ctx context.Context, input GetManyInput) ([]model.ScheduledVip, error)
	GetManyByChannelID(ctx context.Context, channelID string) ([]model.ScheduledVip, error)
	Create(ctx context.Context, input CreateInput) error
	Delete(ctx context.Context, id ulid.ULID) error
	GetByUserAndChannelID(ctx context.Context, userID, channelID string) (model.ScheduledVip, error)
	Update(ctx context.Context, id ulid.ULID, input UpdateInput) error
}

type GetManyInput struct {
	Expired *bool
}

type CreateInput struct {
	ChannelID string
	UserID    string
	RemoveAt  *time.Time
}

type UpdateInput struct {
	RemoveAt *time.Time
}
