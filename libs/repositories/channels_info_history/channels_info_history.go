package channels_info_history

import (
	"context"
	"time"

	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/channels_info_history/model"
)

type Repository interface {
	GetMany(ctx context.Context, input GetManyInput) ([]model.ChannelInfoHistory, error)
	Create(ctx context.Context, input CreateInput) error
}

type UniqueBy int

const (
	UniqueByCategory UniqueBy = iota
	UniqueByTitle
)

type GetManyInput struct {
	ChannelID string
	Platform  *platform.Platform
	After     time.Time // Optional
	Limit     int
	UniqueBy  *UniqueBy
}

type CreateInput struct {
	ChannelID string
	Platform  platform.Platform
	Title     string
	Category  string
}
