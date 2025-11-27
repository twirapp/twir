package chat_translations_settings

import (
	"context"
	"errors"
	"time"

	"github.com/twirapp/kv"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/chat_translation"
	"github.com/twirapp/twir/libs/repositories/chat_translation/model"
)

func New(
	repo chat_translation.Repository,
	kv kv.KV,
) *generic_cacher.GenericCacher[model.ChatTranslation] {
	return generic_cacher.New[model.ChatTranslation](
		generic_cacher.Opts[model.ChatTranslation]{
			KV:        kv,
			KeyPrefix: "cache:twir:chat_translation_settings:channel:",
			LoadFn: func(ctx context.Context, key string) (model.ChatTranslation, error) {
				result, err := repo.GetByChannelID(ctx, key)
				if err != nil {
					if errors.Is(err, chat_translation.ErrSettingsNotFound) {
						return model.ChatTranslationNil, nil
					}

					return model.ChatTranslationNil, err
				}

				return result, nil
			},
			Ttl: 24 * time.Hour,
		},
	)
}
