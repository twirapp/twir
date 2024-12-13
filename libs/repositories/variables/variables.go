package variables

import (
	"context"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/twirapp/twir/libs/repositories/variables/model"
)

type Repository interface {
	GetAllByChannelID(ctx context.Context, channelID string) ([]model.CustomVariable, error)
	CountByChannelID(ctx context.Context, channelID string) (int, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.CustomVariable, error)
	Create(ctx context.Context, input CreateInput) (model.CustomVariable, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.CustomVariable, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateInput struct {
	ChannelID   string
	Name        string
	Description null.String
	Type        model.CustomVarType
	EvalValue   string
	Response    string
}

type UpdateInput struct {
	Name        *string
	Description *string
	Type        *model.CustomVarType
	EvalValue   *string
	Response    *string
}
