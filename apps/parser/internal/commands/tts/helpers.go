package tts

import (
	"encoding/json"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/types/types/api/modules"
	"gorm.io/gorm"
)

func getSettings(channelId, userId string) *modules.TTSSettings {
	db := do.MustInvoke[gorm.DB](di.Provider)

	settings := &model.ChannelModulesSettings{}
	query := db.
		Where(`"channelId" = ?`, channelId).
		Where(`"type" = ?`, "tts")

	if userId != "" {
		query = query.Where(`"userId" = ?`, userId)
	} else {
		query = query.Where(`"userId" IS NULL`)
	}

	err := query.First(&settings).Error
	if err != nil {
		return nil
	}

	data := modules.TTSSettings{}
	err = json.Unmarshal(settings.Settings, &data)
	if err != nil {
		return nil
	}

	return &data
}
