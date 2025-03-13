package bots

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/bots/model"
)

type Repository interface {
	GetDefault(ctx context.Context) (model.Bot, error)
}
