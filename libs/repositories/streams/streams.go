package streams

import (
	"context"
	"time"

	"github.com/twirapp/twir/libs/repositories/streams/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.Stream, error)
	GetList(ctx context.Context) ([]model.Stream, error)
	Save(ctx context.Context, input SaveInput) error
	DeleteByChannelID(ctx context.Context, channelID string) error
	Update(ctx context.Context, channelID string, input UpdateInput) error
}

type SaveInput struct {
	ID           string
	UserId       string
	UserLogin    string
	UserName     string
	GameId       string
	GameName     string
	CommunityIds []string
	Type         string
	Title        string
	ViewerCount  int
	StartedAt    time.Time
	Language     string
	ThumbnailUrl string
	TagIds       []string
	Tags         []string
	IsMature     bool
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
