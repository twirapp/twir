package fetcher

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/twirapp/twir/apps/emotes-cacher/internal/emote"
	"github.com/twirapp/twir/libs/entities/platform"
	seventvapi "github.com/twirapp/twir/libs/integrations/seventv/api"
	seventv "github.com/twirapp/twir/libs/integrations/seventv"
)

type SevenTvGlobalResponse struct {
	Emotes []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"emotes"`
}

func GetChannelSevenTvEmotes(
	ctx context.Context,
	channelPlatform platform.Platform,
	channelPlatformID string,
) ([]emote.Emote, error) {
	client := seventv.NewClient("")

	var (
		profile any
		err error
	)

	switch channelPlatform {
	case platform.PlatformKick:
		profile, err = client.GetProfileByKickId(ctx, channelPlatformID)
	default:
		profile, err = client.GetProfileByTwitchId(ctx, channelPlatformID)
	}
	if err != nil || profile == nil {
		return nil, err
	}

	findActiveSet := func(activeSetID string, sets []seventvapi.TwirSeventvUserEmoteSetsEmoteSet) []emote.Emote {
		for _, set := range sets {
			if set.Id != activeSetID {
				continue
			}

			result := make([]emote.Emote, 0, len(set.Emotes.Items))
			for _, item := range set.Emotes.Items {
				name := item.Alias
				if name == "" {
					name = item.Emote.DefaultName
				}
				result = append(result, emote.Emote{ID: emote.ID(item.Emote.Id), Name: name})
			}

			return result
		}

		return nil
	}

	switch p := profile.(type) {
	case *seventvapi.GetProfileByTwitchIdResponse:
		if p.Users.UserByConnection == nil || p.Users.UserByConnection.Style.ActiveEmoteSet == nil {
			return nil, nil
		}
		return findActiveSet(p.Users.UserByConnection.Style.ActiveEmoteSet.Id, p.Users.UserByConnection.EmoteSets), nil
	case *seventvapi.GetProfileByKickIdResponse:
		if p.Users.UserByConnection == nil || p.Users.UserByConnection.Style.ActiveEmoteSet == nil {
			return nil, nil
		}
		return findActiveSet(p.Users.UserByConnection.Style.ActiveEmoteSet.Id, p.Users.UserByConnection.EmoteSets), nil
	default:
		return nil, nil
	}
}

func GetGlobalSevenTvEmotes(ctx context.Context) ([]emote.Emote, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://7tv.io/v3/emote-sets/global", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data SevenTvGlobalResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	result := make([]emote.Emote, 0, len(data.Emotes))
	for _, e := range data.Emotes {
		result = append(
			result,
			emote.Emote{
				ID:   emote.ID(e.ID),
				Name: e.Name,
			},
		)
	}

	return result, nil
}
