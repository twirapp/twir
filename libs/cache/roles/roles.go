package roles

import (
	"context"
	"time"

	"github.com/twirapp/kv"
	"github.com/twirapp/twir/libs/repositories/roles"
	"github.com/twirapp/twir/libs/repositories/roles/model"

	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
)

func New(
	repo roles.Repository,
	kv kv.KV,
) *generic_cacher.GenericCacher[[]model.Role] {
	return generic_cacher.New(
		generic_cacher.Opts[[]model.Role]{
			KV:        kv,
			KeyPrefix: "cache:twir:roles:channel:",
			LoadFn: func(ctx context.Context, key string) ([]model.Role, error) {
				return repo.GetManyByChannelID(ctx, key)
			},
			Ttl: 24 * time.Hour,
		},
	)
}
