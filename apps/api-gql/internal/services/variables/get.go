package variables

import (
	"context"

	dbmodels "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/services/variables/model"
)

func (c *Service) GetByID(ctx context.Context, id string) (model.Variable, error) {
	entity := dbmodels.ChannelsCustomvars{}
	if err := c.gorm.
		WithContext(ctx).
		Where("id = ?", id).
		First(&entity).Error; err != nil {
		return model.Nil, err
	}

	return c.dbToModel(entity), nil
}

func (c *Service) GetAll(ctx context.Context, channelID string) ([]model.Variable, error) {
	var entities []dbmodels.ChannelsCustomvars
	if err := c.gorm.
		WithContext(ctx).
		Where(`"channelId" = ?`, channelID).
		Order("name ASC").
		Find(&entities).Error; err != nil {
		return nil, err
	}

	converted := make([]model.Variable, 0, len(entities))
	for _, entity := range entities {
		converted = append(converted, c.dbToModel(entity))
	}

	return converted, nil
}
