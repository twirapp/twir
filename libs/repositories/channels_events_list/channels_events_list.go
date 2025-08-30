package channelseventslist

import (
	"context"
	"time"

	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
)

type Repository interface {
	Create(ctx context.Context, input CreateInput) error
	CreateMany(ctx context.Context, inputs []CreateInput) error
	CountBy(ctx context.Context, input CountByInput) (int64, error)
}

type CreateInput struct {
	ChannelID string
	UserID    *string
	Type      model.ChannelEventListItemType
	Data      *model.ChannelsEventsListItemData
}

type CountByInput struct {
	ChannelID    *string
	UserID       *string
	Type         *model.ChannelEventListItemType
	CreatedAtGTE *time.Time
	CreatedAtLTE *time.Time
	CreatedAtGT  *time.Time
	CreatedAtLT  *time.Time
	CreatedAtEQ  *time.Time
}
