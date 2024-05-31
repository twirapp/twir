package seventv

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
)

var ErrSevenTvProfileNotFound = errors.New("7tv profile not found")

func GetProfile(ctx context.Context, twitchUserID string) (ProfileResponse, error) {
	var profile ProfileResponse
	resp, err := req.
		SetHeader("Cache-Control", "no-cache").
		SetContext(ctx).
		SetSuccessResult(&profile).
		Get(fmt.Sprintf("https://7tv.io/v3/users/twitch/%s?t=%v", twitchUserID, time.Now().UnixMilli()))
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

func GetProfileBySevenTvID(ctx context.Context, sevenTvId string) (ProfileResponse, error) {
	var profile ProfileResponse
	resp, err := req.
		SetHeader("Cache-Control", "no-cache").
		SetContext(ctx).
		SetSuccessResult(&profile).
		Get(fmt.Sprintf("https://7tv.io/v3/users/%s?t=%v", sevenTvId, time.Now().UnixMilli()))
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
