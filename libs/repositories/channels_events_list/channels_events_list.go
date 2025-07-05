package channelseventslist

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
)

type Repository interface {
	Create(ctx context.Context, input CreateInput) error
	CreateMany(ctx context.Context, inputs []CreateInput) error
}

type CreateInput struct {
	ChannelID string
	UserID    *string
	Type      model.ChannelEventListItemType
	Data      *model.ChannelsEventsListItemData
}
