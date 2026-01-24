package command_role_cooldown

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/commandrolecooldownentity"
)

type Repository interface {
	GetByCommandID(ctx context.Context, commandID uuid.UUID) ([]commandrolecooldownentity.CommandRoleCooldown, error)
	GetByCommandIDs(
		ctx context.Context,
		commandIDs []uuid.UUID,
	) ([]commandrolecooldownentity.CommandRoleCooldown, error)
	Create(ctx context.Context, input CreateInput) (commandrolecooldownentity.CommandRoleCooldown, error)
	CreateMany(ctx context.Context, inputs []CreateInput) error
	DeleteByCommandID(ctx context.Context, commandID uuid.UUID) error
	DeleteByCommandIDAndRoleID(ctx context.Context, commandID, roleID uuid.UUID) error
}

type CreateInput struct {
	CommandID uuid.UUID
	RoleID    uuid.UUID
	Cooldown  int
}
