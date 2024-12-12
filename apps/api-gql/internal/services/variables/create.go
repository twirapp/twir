package variables

import (
	"context"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/services/variables/model"

	dbmodels "github.com/satont/twir/libs/gomodels"
)

type CreateInput struct {
	ChannelID string
	ActorID   string

	Name        string
	Description *string
	Type        model.CustomVarType
	EvalValue   string
	Response    string
}

func (c *Service) Create(ctx context.Context, data CreateInput) (model.Variable, error) {
	// TODO: write repository
	entity := dbmodels.ChannelsCustomvars{
		ID:          uuid.NewString(),
		Name:        data.Name,
		Description: null.StringFromPtr(data.Description),
		Type:        dbmodels.CustomVarType(data.Type),
		EvalValue:   data.EvalValue,
		Response:    data.Response,
		ChannelID:   data.ChannelID,
	}

	if err := c.gorm.
		WithContext(ctx).
		Create(&entity).Error; err != nil {
		return model.Nil, err
	}

	c.logger.Audit(
		"Variable create",
		audit.Fields{
			OldValue:      nil,
			NewValue:      entity,
			ActorID:       &data.ActorID,
			ChannelID:     &data.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelVariable),
			OperationType: audit.OperationCreate,
			ObjectID:      lo.ToPtr(entity.ID),
		},
	)

	return c.dbToModel(entity), nil
}
