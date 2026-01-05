package scheduled_vips

import (
	"context"
	"time"

	"github.com/google/uuid"
	scheduledvipsentity "github.com/twirapp/twir/libs/entities/scheduled_vips"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (scheduledvipsentity.ScheduledVip, error)
	GetMany(ctx context.Context, input GetManyInput) ([]scheduledvipsentity.ScheduledVip, error)
	Create(ctx context.Context, input CreateInput) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByUserAndChannelID(ctx context.Context, userID, channelID string) (scheduledvipsentity.ScheduledVip, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) error
}

type GetManyInput struct {
	ChannelID  *string
	Expired    *bool
	RemoveType *scheduledvipsentity.RemoveType
}

type CreateInput struct {
	ChannelID  string
	UserID     string
	RemoveAt   *time.Time
	RemoveType *scheduledvipsentity.RemoveType
}

type UpdateInput struct {
	RemoveAt   *time.Time
	RemoveType *scheduledvipsentity.RemoveType
}
