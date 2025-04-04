package chat_translations_settings

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/chat_translation"
	"github.com/twirapp/twir/libs/repositories/chat_translation/model"
)

func New(
	repo chat_translation.Repository,
	redis *redis.Client,
) *generic_cacher.GenericCacher[model.ChatTranslation] {
	return generic_cacher.New[model.ChatTranslation](
		generic_cacher.Opts[model.ChatTranslation]{
			Redis:     redis,
			KeyPrefix: "cache:twir:chat_translation_settings:channel:",
			LoadFn: func(ctx context.Context, key string) (model.ChatTranslation, error) {
				return repo.GetByChannelID(ctx, key)
			},
			Ttl: 24 * time.Hour,
		},
	)
}
