package variables

import (
	"context"

	dbmodels "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func (c *Service) GetByID(ctx context.Context, id string) (entity.Variable, error) {
	e := dbmodels.ChannelsCustomvars{}
	if err := c.gorm.
		WithContext(ctx).
		Where("id = ?", id).
		First(&e).Error; err != nil {
		return entity.CustomVarNil, err
	}

	return c.dbToModel(e), nil
}

func (c *Service) GetAll(ctx context.Context, channelID string) ([]entity.Variable, error) {
	var entities []dbmodels.ChannelsCustomvars
	if err := c.gorm.
		WithContext(ctx).
		Where(`"channelId" = ?`, channelID).
		Order("name ASC").
		Find(&entities).Error; err != nil {
		return nil, err
	}

	converted := make([]entity.Variable, 0, len(entities))
	for _, entity := range entities {
		converted = append(converted, c.dbToModel(entity))
	}

	return converted, nil
}
