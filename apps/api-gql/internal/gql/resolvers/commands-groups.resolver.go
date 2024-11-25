package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/satont/twir/libs/utils"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/mappers"
)

// CommandsGroupsCreate is the resolver for the commandsGroupsCreate field
func (r *mutationResolver) CommandsGroupsCreate(ctx context.Context, opts gqlmodel.CommandsGroupsCreateOpts) (bool, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	var createdCount int64
	if err := r.gorm.
		WithContext(ctx).
		Model(&model.ChannelCommandGroup{}).
		Where(`"channelId" = ?`, dashboardId).
		Count(&createdCount).
		Error; err != nil {
		return false, err
	}

	if createdCount >= 10 {
		return false, fmt.Errorf("you can have only 10 command groups")
	}

	entity := model.ChannelCommandGroup{
		ID:        uuid.NewString(),
		ChannelID: dashboardId,
		Name:      opts.Name,
		Color:     opts.Color,
	}

	if err := r.gorm.WithContext(ctx).Create(&entity).Error; err != nil {
		return false, err
	}

	r.logger.Audit(
		"New command group",
		audit.Fields{
			NewValue:      entity,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardId),
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelCommandGroup),
			OperationType: audit.OperationCreate,
			ObjectID:      &entity.ID,
		},
	)

	return true, nil
}

// CommandsGroupsUpdate is the resolver for the commandsGroupsUpdate field.
func (r *mutationResolver) CommandsGroupsUpdate(ctx context.Context, id string, opts gqlmodel.CommandsGroupsUpdateOpts) (bool, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	entity := model.ChannelCommandGroup{}
	if err := r.gorm.
		WithContext(ctx).
		Where(`id = ? AND "channelId" = ?`, id, dashboardId).
		First(&entity).
		Error; err != nil {
		return false, fmt.Errorf("group not found: %w", err)
	}

	var entityCopy model.ChannelCommandGroup
	if err := utils.DeepCopy(entity, &entityCopy); err != nil {
		return false, err
	}

	if opts.Name.IsSet() {
		entity.Name = *opts.Name.Value()
	}

	if opts.Color.IsSet() {
		entity.Color = *opts.Color.Value()
	}

	if err := r.gorm.WithContext(ctx).Save(&entity).Error; err != nil {
		return false, err
	}

	if err := r.cachedCommandsClient.Invalidate(ctx, dashboardId); err != nil {
		r.logger.Error("failed to invalidate commands cache", slog.Any("err", err))
	}

	r.logger.Audit(
		"Command group update",
		audit.Fields{
			OldValue:      entityCopy,
			NewValue:      entity,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardId),
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelCommandGroup),
			OperationType: audit.OperationUpdate,
			ObjectID:      &entity.ID,
		},
	)

	return true, nil
}

// CommandsGroupsRemove is the resolver for the commandsGroupsRemove field.
func (r *mutationResolver) CommandsGroupsRemove(ctx context.Context, id string) (bool, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	entity := model.ChannelCommandGroup{}
	if err := r.gorm.
		WithContext(ctx).
		Where(`id = ? AND "channelId" = ?`, id, dashboardId).
		First(&entity).
		Error; err != nil {
		return false, fmt.Errorf("group not found: %w", err)
	}

	if err := r.gorm.
		WithContext(ctx).
		Delete(&entity).
		Error; err != nil {
		return false, err
	}

	if err := r.cachedCommandsClient.Invalidate(ctx, dashboardId); err != nil {
		r.logger.Error("failed to invalidate commands cache", slog.Any("err", err))
	}

	r.logger.Audit(
		"Command group remove",
		audit.Fields{
			OldValue:      entity,
			NewValue:      nil,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardId),
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelCommandGroup),
			OperationType: audit.OperationDelete,
			ObjectID:      &entity.ID,
		},
	)

	return true, nil
}

// CommandsGroups is the resolver for the commandsGroups field.
func (r *queryResolver) CommandsGroups(ctx context.Context) ([]gqlmodel.CommandGroup, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	var entities []model.ChannelCommandGroup
	if err := r.gorm.
		WithContext(ctx).
		Where(`"channelId" = ?`, dashboardId).
		Find(&entities).
		Error; err != nil {
		return nil, err
	}

	var result []gqlmodel.CommandGroup
	for _, entity := range entities {
		result = append(
			result, gqlmodel.CommandGroup{
				ID:    entity.ID,
				Name:  entity.Name,
				Color: entity.Color,
			},
		)
	}

	return result, nil
}
