package spotify

import (
	"context"
	"fmt"
	"strings"

	"github.com/imroc/req/v3"
)

type recentTracksResponse struct {
	Href    string `json:"href"`
	Limit   int    `json:"limit"`
	Next    string `json:"next"`
	Cursors struct {
		After  string `json:"after"`
		Before string `json:"before"`
	} `json:"cursors"`
	Total int `json:"total"`
	Items []struct {
		Track struct {
			Album struct {
				AlbumType        string   `json:"album_type"`
				TotalTracks      int      `json:"total_tracks"`
				AvailableMarkets []string `json:"available_markets"`
				ExternalUrls     struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href   string `json:"href"`
				Id     string `json:"id"`
				Images []struct {
					Url    string `json:"url"`
					Height int    `json:"height"`
					Width  int    `json:"width"`
				} `json:"images"`
				Name                 string `json:"name"`
				ReleaseDate          string `json:"release_date"`
				ReleaseDatePrecision string `json:"release_date_precision"`
				Restrictions         struct {
					Reason string `json:"reason"`
				} `json:"restrictions"`
				Type    string `json:"type"`
				Uri     string `json:"uri"`
				Artists []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href string `json:"href"`
					Id   string `json:"id"`
					Name string `json:"name"`
					Type string `json:"type"`
					Uri  string `json:"uri"`
				} `json:"artists"`
			} `json:"album"`
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				Id   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				Uri  string `json:"uri"`
			} `json:"artists"`
			AvailableMarkets []string `json:"available_markets"`
			DiscNumber       int      `json:"disc_number"`
			DurationMs       int      `json:"duration_ms"`
			Explicit         bool     `json:"explicit"`
			ExternalIds      struct {
				Isrc string `json:"isrc"`
				Ean  string `json:"ean"`
				Upc  string `json:"upc"`
			} `json:"external_ids"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href       string `json:"href"`
			Id         string `json:"id"`
			IsPlayable bool   `json:"is_playable"`
			LinkedFrom struct {
			} `json:"linked_from"`
			Restrictions struct {
				Reason string `json:"reason"`
			} `json:"restrictions"`
			Name        string `json:"name"`
			Popularity  int    `json:"popularity"`
			PreviewUrl  string `json:"preview_url"`
			TrackNumber int    `json:"track_number"`
			Type        string `json:"type"`
			Uri         string `json:"uri"`
			IsLocal     bool   `json:"is_local"`
		} `json:"track"`
		PlayedAt string `json:"played_at"`
		Context  struct {
			Type         string `json:"type"`
			Href         string `json:"href"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Uri string `json:"uri"`
		} `json:"context"`
	} `json:"items"`
}

type RecentTrack struct {
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Image    string `json:"image"`
	PlayedAt string `json:"playedAt"`
}

type GetRecentTracksInput struct {
	Limit int
}

var ErrNoNeededScope = fmt.Errorf("no needed scope")

func (c *Spotify) GetRecentTracks(ctx context.Context, input GetRecentTracksInput) (
	[]RecentTrack,
	error,
) {
	limit := input.Limit
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	data := recentTracksResponse{}
	resp, err := req.R().
		SetBearerAuthToken(c.channelIntegration.AccessToken).
		SetSuccessResult(&data).
		SetQueryParams(
			map[string]string{
				"limit": fmt.Sprintf("%d", limit),
			},
		).
		Get("https://api.spotify.com/v1/me/player/recently-played")
	if resp != nil && resp.StatusCode == 401 && !c.isRetry {
		c.isRetry = true
		c.refreshToken(ctx)
		return c.GetRecentTracks(ctx, input)
	}

	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("cannot get recent tracks: %s", "empty response")
	}

	if !resp.IsSuccessState() {
		if resp.StatusCode == 403 {
			return nil, ErrNoNeededScope
		}

		return nil, fmt.Errorf("cannot get recent tracks: %s", resp.String())
	}

	recentTracks := make([]RecentTrack, 0, len(data.Items))
	for _, item := range data.Items {
		var image string
		if len(item.Track.Album.Images) > 0 {
			image = item.Track.Album.Images[0].Url
		}
		artists := make([]string, len(item.Track.Artists))
		for i, artist := range item.Track.Artists {
			artists[i] = artist.Name
		}

		recentTracks = append(
			recentTracks,
			RecentTrack{
				Title:    item.Track.Name,
				Artist:   strings.Join(artists, ", "),
				Image:    image,
				PlayedAt: item.PlayedAt,
			},
		)
	}

	return recentTracks, nil
}
