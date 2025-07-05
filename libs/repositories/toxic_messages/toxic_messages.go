package toxic_messages

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/toxic_messages/model"
)

type Repository interface {
	Create(ctx context.Context, input CreateInput) error
	GetList(ctx context.Context, input GetListInput) (GetListOutput, error)
}

type CreateInput struct {
	ChannelID     string
	ReplyToUserID *string
	Text          string
}

type GetListInput struct {
	Page    int
	PerPage int
}

type GetListOutput struct {
	Items []model.ToxicMessage
	Total int
}
