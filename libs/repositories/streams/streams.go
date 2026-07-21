package streams

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/streams/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID uuid.UUID, platform platform.Platform) (model.Stream, error)
	GetByUserID(ctx context.Context, userID string, platform platform.Platform) (model.Stream, error)
	GetListByChannelID(ctx context.Context, channelID uuid.UUID) ([]model.Stream, error)
	GetList(ctx context.Context) ([]model.Stream, error)
	Count(ctx context.Context) (uint64, error)
	Save(ctx context.Context, input SaveInput) error
	DeleteByChannelID(ctx context.Context, channelID uuid.UUID, platform platform.Platform) error
	Update(ctx context.Context, channelID uuid.UUID, platform platform.Platform, input UpdateInput) error
}

type SaveInput struct {
	ID           string
	ChannelID    uuid.UUID
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
	Platform     platform.Platform
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
