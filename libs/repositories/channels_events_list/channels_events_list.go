package channelseventslist

import (
	"context"
	"errors"
	"time"

	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
)

var ErrEmptyPlatform = errors.New("channels_events_list: platform is required")

type Repository interface {
	Create(ctx context.Context, input CreateInput) error
	CreateMany(ctx context.Context, inputs []CreateInput) error
	CountBy(ctx context.Context, input CountByInput) (int64, error)
}

type CreateInput struct {
	ChannelID string
	UserID    *string
	Platform  platform.Platform
	Type      model.ChannelEventListItemType
	Data      *model.ChannelsEventsListItemData
}

type CountByInput struct {
	ChannelID    *string
	UserID       *string
	Platform     *platform.Platform
	Type         *model.ChannelEventListItemType
	CreatedAtGTE *time.Time
	CreatedAtLTE *time.Time
	CreatedAtGT  *time.Time
	CreatedAtLT  *time.Time
	CreatedAtEQ  *time.Time
}
