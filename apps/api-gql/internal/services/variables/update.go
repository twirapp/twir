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
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

type UpdateInput struct {
	ID        string
	ChannelID string
	ActorID   string

	Name        *string
	Description *string
	Type        *entity.CustomVarType
	EvalValue   *string
	Response    *string
}

func (c *Service) Update(ctx context.Context, data UpdateInput) (entity.Variable, error) {
	e := dbmodels.ChannelsCustomvars{}
	if err := c.gorm.
		WithContext(ctx).
		Where(`"channelId" = ? AND id = ?`, data.ChannelID, data.ID).
		First(&e).Error; err != nil {
		return entity.CustomVarNil, err
	}

	var entityCopy dbmodels.ChannelsCustomvars
	if err := utils.DeepCopy(&e, &entityCopy); err != nil {
		return entity.CustomVarNil, err
	}

	if data.Name != nil {
		e.Name = *data.Name
	}

	if data.Description != nil {
		e.Description = null.StringFromPtr(data.Description)
	}

	if data.Type != nil {
		e.Type = dbmodels.CustomVarType(*data.Type)
	}

	if data.EvalValue != nil {
		e.EvalValue = *data.EvalValue
	}

	if data.Response != nil {
		e.Response = *data.Response
	}

	if err := c.gorm.
		WithContext(ctx).
		Save(&e).Error; err != nil {
		return entity.CustomVarNil, err
	}

	c.logger.Audit(
		"Variable update",
		audit.Fields{
			OldValue:      entityCopy,
			NewValue:      e,
			ActorID:       &data.ActorID,
			ChannelID:     &data.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelVariable),
			OperationType: audit.OperationUpdate,
			ObjectID:      lo.ToPtr(e.ID),
		},
	)

	return c.dbToModel(e), nil
}
