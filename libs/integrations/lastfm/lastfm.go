package lastfm

import (
	"errors"
	"fmt"

	api "github.com/shkh/lastfm-go/lastfm"
)

type Opts struct {
	ApiKey       string
	ClientSecret string
	SessionKey   string
}

type Lastfm struct {
	api  *api.Api
	user *api.UserGetInfo
}

func New(opts Opts) (*Lastfm, error) {
	if opts.ApiKey == "" || opts.ClientSecret == "" || opts.SessionKey == "" {
		return nil, errors.New("lastfm api key, client secret, and session key are required")
	}

	lfm := &Lastfm{
		api: api.New(opts.ApiKey, opts.ClientSecret),
	}

	lfm.api.SetSession(opts.SessionKey)

	user, err := lfm.api.User.GetInfo(map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	lfm.user = &user

	return lfm, nil
}

type Track struct {
	Title     string
	Artist    string
	Image     string
	PlayedUTS string
}

func (c *Lastfm) GetTrack() (*Track, error) {
	tracks, err := c.api.User.GetRecentTracks(
		map[string]interface{}{
			"limit": "1",
			"user":  c.user.Name,
		},
	)
	if err != nil {
		return nil, err
	}

	if len(tracks.Tracks) == 0 || tracks.Tracks[0].NowPlaying != "true" {
		return nil, nil
	}

	track := tracks.Tracks[0]
	var cover string
	if len(track.Images) > 0 {
		cover = track.Images[0].Url
	}

	return &Track{
		Title:     track.Name,
		Artist:    track.Artist.Name,
		Image:     cover,
		PlayedUTS: track.Date.Date,
	}, nil
}

func (c *Lastfm) GetRecentTracks(limit int) ([]Track, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	tracks, err := c.api.User.GetRecentTracks(
		map[string]interface{}{
			"limit": fmt.Sprintf("%d", limit),
			"user":  c.user.Name,
		},
	)
	if err != nil {
		return nil, err
	}

	recentTracks := make([]Track, 0, len(tracks.Tracks))
	for _, track := range tracks.Tracks {
		if track.NowPlaying == "true" {
			continue
		}

		var cover string
		if len(track.Images) > 0 {
			cover = track.Images[0].Url
		}

		recentTracks = append(
			recentTracks,
			Track{
				Title:     track.Name,
				Artist:    track.Artist.Name,
				Image:     cover,
				PlayedUTS: track.Date.Uts,
			},
		)
	}

	return recentTracks, nil
}
