package tts

import (
	"context"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/types/types/api/modules"
	"gorm.io/gorm"

	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
)

func NewTTSSettings(
	gorm *gorm.DB,
	redis *redis.Client,
) *generic_cacher.GenericCacher[modules.TTSSettings] {
	return generic_cacher.New[modules.TTSSettings](
		generic_cacher.Opts[modules.TTSSettings]{
			Redis:     redis,
			KeyPrefix: "cache:twir:tts-settings:channel:",
			LoadFn: func(ctx context.Context, key string) (modules.TTSSettings, error) {
				entity := &model.ChannelModulesSettings{}
				err := gorm.WithContext(ctx).
					Where(`"channelId" = ?`, key).
					Where(`"userId" IS NULL`).
					Where(`"type" = ?`, "tts").
					First(entity).
					Error
				if err != nil {
					return modules.TTSSettings{}, err
				}

				data := modules.TTSSettings{}
				err = json.Unmarshal(entity.Settings, &data)
				if err != nil {
					return modules.TTSSettings{}, err
				}

				return data, nil
			},
			Ttl: 24 * time.Hour,
		},
	)
}
