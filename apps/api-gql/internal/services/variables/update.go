package variables

import (
	"context"

	"github.com/guregu/null"
	"github.com/samber/lo"
	dbmodels "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/satont/twir/libs/utils"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/services/variables/model"
)

type UpdateInput struct {
	ID        string
	ChannelID string
	ActorID   string

	Name        *string
	Description *string
	Type        *model.CustomVarType
	EvalValue   *string
	Response    *string
}

func (c *Service) Update(ctx context.Context, data UpdateInput) (model.Variable, error) {
	entity := dbmodels.ChannelsCustomvars{}
	if err := c.gorm.
		WithContext(ctx).
		Where(`"channelId" = ? AND id = ?`, data.ChannelID, data.ID).
		First(&entity).Error; err != nil {
		return model.Nil, err
	}

	var entityCopy dbmodels.ChannelsCustomvars
	if err := utils.DeepCopy(&entity, &entityCopy); err != nil {
		return model.Nil, err
	}

	if data.Name != nil {
		entity.Name = *data.Name
	}

	if data.Description != nil {
		entity.Description = null.StringFromPtr(data.Description)
	}

	if data.Type != nil {
		entity.Type = dbmodels.CustomVarType(*data.Type)
	}

	if data.EvalValue != nil {
		entity.EvalValue = *data.EvalValue
	}

	if data.Response != nil {
		entity.Response = *data.Response
	}

	if err := c.gorm.
		WithContext(ctx).
		Save(&entity).Error; err != nil {
		return model.Nil, err
	}

	c.logger.Audit(
		"Variable update",
		audit.Fields{
			OldValue:      entityCopy,
			NewValue:      entity,
			ActorID:       &data.ActorID,
			ChannelID:     &data.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelVariable),
			OperationType: audit.OperationUpdate,
			ObjectID:      lo.ToPtr(entity.ID),
		},
	)

	return c.dbToModel(entity), nil
}
