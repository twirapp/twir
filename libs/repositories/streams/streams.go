package streams

import (
	"context"
	"time"

	"github.com/twirapp/twir/libs/repositories/streams/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.Stream, error)
	GetList(ctx context.Context) ([]model.Stream, error)
	Update(ctx context.Context, channelID string, input UpdateInput) error
}

type UpdateInput struct {
	GameId       *string
	GameName     *string
	CommunityIds []string
	Type         *string
	Title        *string
	ViewerCount  *int
	StartedAt    *time.Time
	Language     *string
	ThumbnailUrl *string
	TagIds       []string
	Tags         []string
	IsMature     *bool
}
