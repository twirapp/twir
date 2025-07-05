package chatalerts

import (
	"context"
	"errors"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	model "github.com/satont/twir/libs/gomodels"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"gorm.io/gorm"
)

type ChatAlert struct {
	model.ChannelModulesSettings
	ParsedSettings model.ChatAlertsSettings
}

var ErrChatAlertNotFound = errors.New("not found")

func New(
	db *gorm.DB,
	redis *redis.Client,
) *generic_cacher.GenericCacher[ChatAlert] {
	return generic_cacher.New[ChatAlert](
		generic_cacher.Opts[ChatAlert]{
			Redis:     redis,
			KeyPrefix: "cache:twir:chat_alerts:channel:",
			LoadFn: func(ctx context.Context, key string) (ChatAlert, error) {
				entity := model.ChannelModulesSettings{}
				if err := db.
					WithContext(ctx).
					Where(
						`"channelId" = ? AND "userId" IS NULL AND type = 'chat_alerts'`,
						key,
					).First(&entity).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return ChatAlert{}, ErrChatAlertNotFound
					}

					return ChatAlert{}, err
				}

				parsedSettings := model.ChatAlertsSettings{}
				if err := json.Unmarshal(entity.Settings, &parsedSettings); err != nil {
					return ChatAlert{}, err
				}

				return ChatAlert{
					ChannelModulesSettings: entity,
					ParsedSettings:         parsedSettings,
				}, nil
			},
			Ttl: 24 * time.Hour,
		},
	)
}
