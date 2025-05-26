package channelseventslist

import (
	"context"

	model "github.com/satont/twir/libs/gomodels"
)

type Repository interface {
	CreateMany(ctx context.Context, inputs []CreateInput) error
}

type CreateInput struct {
	ChannelID string
	UserID    *string
	Type      model.ChannelEventListItemType
	Data      *model.ChannelsEventsListItemData
}
