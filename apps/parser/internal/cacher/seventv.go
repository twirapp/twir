package cacher

import (
	"context"

	seventvintegration "github.com/twirapp/twir/libs/integrations/seventv"
	seventvintegrationapi "github.com/twirapp/twir/libs/integrations/seventv/api"
)

func (c *cacher) GetSeventvProfileGetTwitchId(ctx context.Context, userId string) (
	*seventvintegrationapi.TwirSeventvUser,
	error,
) {
	if c.cache.seventvprofile != nil {
		return c.cache.seventvprofile, nil
	}

	client := seventvintegration.NewClient(c.services.Config.SevenTvToken)

	profile, err := client.GetProfileByTwitchId(ctx, userId)
	if err != nil {
		return nil, err
	}

	c.cache.seventvprofile = &profile.Users.UserByConnection.TwirSeventvUser

	return c.cache.seventvprofile, nil
}
