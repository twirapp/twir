package tts

import (
	"context"

	"github.com/goccy/go-json"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/types/types/api/modules"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Gorm *gorm.DB
}

func New(opts Opts) *Service {
	return &Service{
		gorm: opts.Gorm,
	}
}

type Service struct {
	gorm *gorm.DB
}

func (s *Service) mapToEntity(
	userId string,
	isChannelOwner bool,
	data []byte,
) entity.TTSUserSettings {
	settings := modules.TTSSettings{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return entity.TTSUserSettings{}
	}

	return entity.TTSUserSettings{
		UserID:         userId,
		Rate:           settings.Rate,
		Pitch:          settings.Pitch,
		Volume:         settings.Volume,
		Voice:          settings.Voice,
		IsChannelOwner: isChannelOwner,
	}
}

func (s *Service) GetTTSUsersSettings(
	ctx context.Context,
	channelID string,
) ([]entity.TTSUserSettings, error) {
	var entities []model.ChannelModulesSettings
	if err := s.gorm.WithContext(ctx).
		Where(
			`"channelId" = ? AND "userId" IS NOT NULL AND type = ?`,
			channelID,
			"tts",
		).
		Order(`"userId" desc`).
		Find(&entities).Error; err != nil {
		return nil, err
	}

	mappedEntities := make([]entity.TTSUserSettings, 0, len(entities))

	channelSettings := model.ChannelModulesSettings{}
	if err := s.gorm.WithContext(ctx).Where(
		`"channelId" = ? AND "userId" IS NULL AND type = ?`,
		channelID,
		"tts",
	).Find(&channelSettings).Error; err != nil {
		return nil, err
	}

	if channelSettings.ID != "" {
		mappedEntities = append(
			mappedEntities,
			s.mapToEntity(channelID, true, channelSettings.Settings),
		)
	}

	for _, e := range entities {
		mappedEntities = append(
			mappedEntities,
			s.mapToEntity(e.UserId.String, false, e.Settings),
		)
	}

	return mappedEntities, nil
}
