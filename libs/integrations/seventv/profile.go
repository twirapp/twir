package seventv

import (
	"context"
	"errors"
	"fmt"

	"github.com/imroc/req/v3"
)

var ErrSevenTvProfileNotFound = errors.New("7tv profile not found")

func GetProfile(ctx context.Context, twitchUserID string) (SevenTvProfileResponse, error) {
	var profile SevenTvProfileResponse
	resp, err := req.
		SetContext(ctx).
		SetSuccessResult(&profile).
		Get("https://7tv.io/v3/users/twitch/" + twitchUserID)
	if err != nil {
		return profile, err
	}

	if !resp.IsSuccessState() {
		if resp.StatusCode == 404 {
			return profile, ErrSevenTvProfileNotFound
		}
		return profile, fmt.Errorf("failed to get 7tv data: %s", resp.String())
	}

	return profile, nil
}
