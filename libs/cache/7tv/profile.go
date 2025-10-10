package seventv

import (
	"context"
	"time"

	"github.com/twirapp/kv"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	config "github.com/twirapp/twir/libs/config"
	seventvintegration "github.com/twirapp/twir/libs/integrations/seventv"
	seventvintegrationapi "github.com/twirapp/twir/libs/integrations/seventv/api"
)

func New(
	kv kv.KV,
	cfg config.Config,
) *generic_cacher.GenericCacher[seventvintegrationapi.TwirSeventvUser] {
	return generic_cacher.New[seventvintegrationapi.TwirSeventvUser](
		generic_cacher.Opts[seventvintegrationapi.TwirSeventvUser]{
			KV:        kv,
			KeyPrefix: "cache:twir:seventv:profilev2:",
			LoadFn: func(ctx context.Context, key string) (
				seventvintegrationapi.TwirSeventvUser,
				error,
			) {
				client := seventvintegration.NewClient(cfg.SevenTvToken)

				profile, err := client.GetProfileByTwitchId(ctx, key)
				if err != nil {
					return seventvintegrationapi.TwirSeventvUser{}, err
				}

				return profile.Users.UserByConnection.TwirSeventvUser, nil
			},
			Ttl: 5 * time.Minute,
		},
	)
}
