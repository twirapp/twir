package tts

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/types/types/api/modules"
	"github.com/twirapp/twir/libs/api/messages/tts_unprotected"
	"github.com/twitchtv/twirp"
)

type Tts struct {
	*impl_deps.Deps
}

func (c *Tts) GetTTSChannelSettings(
	ctx context.Context,
	req *tts_unprotected.GetChannelSettingsRequest,
) (*tts_unprotected.Settings, error) {
	channel := &model.Channels{}
	if err := c.Db.
		WithContext(ctx).
		Where(`id = ?`, req.ChannelId).
		Joins("User").
		First(channel).Error; err != nil {
		return nil, err
	}

	if channel.User.IsBanned {
		return &tts_unprotected.Settings{}, nil
	}

	entity := model.ChannelModulesSettings{}
	if err := c.Db.WithContext(ctx).Where(
		`"channelId" = ? AND "userId" IS NULL AND type = ?`,
		req.ChannelId,
		"tts",
	).Find(&entity).Error; err != nil {
		return nil, err
	}

	if entity.ID == "" {
		return nil, twirp.NotFoundError("settings not found")
	}

	settings := modules.TTSSettings{}
	if err := json.Unmarshal(entity.Settings, &settings); err != nil {
		return nil, err
	}

	return &tts_unprotected.Settings{
		UserId: "",
		Rate:   uint32(settings.Rate),
		Volume: uint32(settings.Volume),
		Pitch:  uint32(settings.Pitch),
		Voice:  settings.Voice,
	}, nil
}

func (c *Tts) GetTTSUsersSettings(
	ctx context.Context,
	req *tts_unprotected.GetUsersSettingsRequest,
) (*tts_unprotected.GetUsersSettingsResponse, error) {
	var entities []model.ChannelModulesSettings
	if err := c.Db.WithContext(ctx).Where(
		`"channelId" = ? AND "userId" IS NOT NULL AND type = ?`,
		req.ChannelId,
		"tts",
	).Find(&entities).Error; err != nil {
		return nil, err
	}

	resultSettings := make([]*tts_unprotected.Settings, 0, len(entities))
	for _, e := range entities {
		s := modules.TTSSettings{}
		if err := json.Unmarshal(e.Settings, &s); err != nil {
			return nil, err
		}
		resultSettings = append(
			resultSettings, &tts_unprotected.Settings{
				UserId: e.UserId.String,
				Rate:   uint32(s.Rate),
				Volume: uint32(s.Volume),
				Pitch:  uint32(s.Pitch),
				Voice:  s.Voice,
			},
		)
	}

	return &tts_unprotected.GetUsersSettingsResponse{
		Settings: resultSettings,
	}, nil
}
