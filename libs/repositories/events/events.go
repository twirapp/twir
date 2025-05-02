package events

import (
	"context"
	"errors"

	"github.com/twirapp/twir/libs/repositories/events/model"
)

type Repository interface {
	GetManyByChannelID(ctx context.Context, channelID string) ([]model.Event, error)
	GetByID(ctx context.Context, id string) (model.Event, error)
	Create(ctx context.Context, input CreateInput) (model.Event, error)
	Update(ctx context.Context, id string, input UpdateInput) (model.Event, error)
	Delete(ctx context.Context, id string) error
}

type CreateInput struct {
	ChannelID   string
	Type        string
	RewardID    *string
	CommandID   *string
	KeywordID   *string
	Description string
	Enabled     bool
	OnlineOnly  bool
	Operations  []OperationInput
}

type UpdateInput struct {
	Type        *string
	RewardID    *string
	CommandID   *string
	KeywordID   *string
	Description *string
	Enabled     *bool
	OnlineOnly  *bool
	Operations  *[]OperationInput
}

type OperationInput struct {
	Type           string
	Input          *string
	Delay          int
	Repeat         int
	UseAnnounce    bool
	TimeoutTime    int
	TimeoutMessage *string
	Target         *string
	Enabled        bool
	Filters        []OperationFilterInput
}

type OperationFilterInput struct {
	Type  string
	Left  string
	Right string
}

var ErrNotFound = errors.New("event not found")
