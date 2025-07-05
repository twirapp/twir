package commands

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/commands/model"
)

type Repository interface {
	GetManyByChannelID(ctx context.Context, channelID string) ([]model.Command, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Command, error)
	Create(ctx context.Context, input CreateInput) (model.Command, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.Command, error)
}

type CreateInput struct {
	ChannelID                 string
	Name                      string
	Cooldown                  int
	CooldownType              string
	Enabled                   bool
	Aliases                   []string
	Description               string
	Visible                   bool
	IsReply                   bool
	KeepResponsesOrder        bool
	DeniedUsersIDS            []string
	AllowedUsersIDS           []string
	RolesIDS                  []string
	OnlineOnly                bool
	OfflineOnly               bool
	CooldownRolesIDs          []string
	EnabledCategories         []string
	RequiredWatchTime         int
	RequiredMessages          int
	RequiredUsedChannelPoints int
	GroupID                   *uuid.UUID
	ExpiresAt                 *time.Time
	ExpiresType               *string
}

type UpdateInput struct {
	Name                      *string
	Cooldown                  *int
	CooldownType              *string
	Enabled                   *bool
	Aliases                   []string
	Description               *string
	Visible                   *bool
	IsReply                   *bool
	KeepResponsesOrder        *bool
	DeniedUsersIDS            []string
	AllowedUsersIDS           []string
	RolesIDS                  []string
	OnlineOnly                *bool
	OfflineOnly               *bool
	CooldownRolesIDs          []string
	EnabledCategories         []string
	RequiredWatchTime         *int
	RequiredMessages          *int
	RequiredUsedChannelPoints *int
	GroupID                   *uuid.UUID
	ExpiresAt                 *time.Time
	ExpiresType               *string
}
