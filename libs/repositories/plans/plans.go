package plans

import (
	"context"

	"github.com/twirapp/twir/libs/entities/plan"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (plan.Plan, error)
	GetByNameID(ctx context.Context, nameID string) (plan.Plan, error)
	GetByChannelID(ctx context.Context, channelID string) (plan.Plan, error)
}
