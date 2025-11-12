package usersstats

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/users_stats/model"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.UserStat, error)
	GetByUserAndChannelID(ctx context.Context, userID, channelID string) (*model.UserStat, error)
	Create(ctx context.Context, input CreateInput) (*model.UserStat, error)
	CreateOrUpdate(ctx context.Context, userID, channelID string, input UpdateInput) (
		*model.UserStat,
		error,
	)
}

type CreateInput struct {
	UserID            string
	ChannelID         string
	Messages          int32
	Watched           int64
	UsedChannelPoints int64
	IsMod             bool
	IsVip             bool
	IsSubscriber      bool
	Reputation        int64
	Emotes            int
}

type UpdateInput struct {
	NumberFields UpdateNumberFieldsInput
	IsMod        *bool
	IsVip        *bool
	IsSubscriber *bool
}

type IncrementInputFieldName string

func (c IncrementInputFieldName) String() string {
	return string(c)
}

const (
	IncrementInputFieldMessages          IncrementInputFieldName = "messages"
	IncrementInputFieldWatched           IncrementInputFieldName = "watched"
	IncrementInputFieldUsedChannelPoints IncrementInputFieldName = "usedChannelPoints"
	IncrementInputFieldReputation        IncrementInputFieldName = "reputation"
	IncrementInputFieldEmotes            IncrementInputFieldName = "emotes"
)

type NumberFieldUpdate struct {
	Count       int
	IsIncrement bool
}

type UpdateNumberFieldsInput map[IncrementInputFieldName]NumberFieldUpdate
