package variables

import (
	"context"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/samber/lo"
	dbmodels "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

type CreateInput struct {
	ChannelID string
	ActorID   string

	Name        string
	Description *string
	Type        entity.CustomVarType
	EvalValue   string
	Response    string
}

func (c *Service) Create(ctx context.Context, data CreateInput) (entity.CustomVariable, error) {
	// TODO: write repository
	e := dbmodels.ChannelsCustomvars{
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
		Create(&e).Error; err != nil {
		return entity.CustomVarNil, err
	}

	c.logger.Audit(
		"Variable create",
		audit.Fields{
			OldValue:      nil,
			NewValue:      e,
			ActorID:       &data.ActorID,
			ChannelID:     &data.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelVariable),
			OperationType: audit.OperationCreate,
			ObjectID:      lo.ToPtr(e.ID),
		},
	)

	return c.dbToModel(e), nil
}
