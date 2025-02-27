package spotify

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
)

type spotifyCurrentPlayingArtist struct {
	Name string `json:"name"`
}

type spotifyCurrentPlayingImage struct {
	URL string `json:"url"`
}

type spotifyCurrentPlayingAlbum struct {
	Images []spotifyCurrentPlayingImage `json:"images"`
}

type spotifyCurrentPlayingTrack struct {
	Artists []spotifyCurrentPlayingArtist `json:"artists"`
	Name    string                        `json:"name"`
	Album   spotifyCurrentPlayingAlbum    `json:"album"`
}

type spotifyCurrentPlayingResponse struct {
	Track     *spotifyCurrentPlayingTrack `json:"item"`
	IsPlaying bool                        `json:"is_playing"`
}

func (c *Spotify) GetTrack(ctx context.Context) (*GetTrackResponse, error) {
	if !slices.Contains(c.channelIntegration.Scopes, "user-read-playback-state") {
		return c.getTrackByCurrentPlayingTrack(ctx)
	}

	return c.getTrackByPlayerState(ctx)
}

type GetTrackResponsePlaylistMeta struct {
	Name      string   `json:"name"`
	Images    []string `json:"images"`
	Followers int      `json:"followers"`
}

type GetTrackResponsePlaylist struct {
	Meta        *GetTrackResponsePlaylistMeta `json:"meta"`        // exists only if the playlist is public
	Uri         string                        `json:"uri"`         // for example spotify:playlist:37i9dQZF1E8O7Yz282BRuZ
	ExternalUrl string                        `json:"externalUrl"` // for example https://open.spotify.com/playlist/37i9dQZF1E8O7Yz282BRuZ
}

type GetTrackResponse struct {
	Playlist  *GetTrackResponsePlaylist `json:"playlist"` // exists only if the track is from a playlist and integration has the required scopes
	Title     string                    `json:"title"`
	Artist    string                    `json:"artist"`
	Image     string                    `json:"image"`
	IsPlaying bool                      `json:"isPlaying"`
}

func (c *Spotify) getTrackByCurrentPlayingTrack(ctx context.Context) (*GetTrackResponse, error) {
	data := spotifyCurrentPlayingResponse{}
	resp, err := req.R().
		SetBearerAuthToken(c.channelIntegration.AccessToken).
		SetSuccessResult(&data).
		Get("https://api.spotify.com/v1/me/player/currently-playing")

	if resp.StatusCode == 401 && !c.isRetry {
		c.isRetry = true
		defer func() {
			c.isRetry = false
		}()
		if err := c.refreshToken(ctx); err != nil {
			return nil, err
		}

		return c.GetTrack(ctx)
	}
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("cannot get spotify track: %s", resp.String())
	}

	if data.Track == nil {
		return nil, nil
	}

	artistsMap := lo.Map(
		data.Track.Artists, func(artist spotifyCurrentPlayingArtist, _ int) string {
			return artist.Name
		},
	)

	var imageUrl string
	if len(data.Track.Album.Images) > 0 {
		imageUrl = data.Track.Album.Images[0].URL
	}

	return &GetTrackResponse{
		Artist:    strings.Join(artistsMap, ", "),
		Title:     data.Track.Name,
		Image:     imageUrl,
		IsPlaying: data.IsPlaying,
	}, nil
}

type spotifyPlayerStateResponse struct {
	Context struct {
		Type         string `json:"type"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		URI string `json:"uri"`
	} `json:"context"`

	Item struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Album struct {
			Name   string `json:"name"`
			Images []struct {
				URL string `json:"url"`
			}
		}
		Artists []struct {
			Name         string `json:"name"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			}
			Uri string `json:"uri"`
		}
		TrackNumber int `json:"track_number"`
		Popularity  int `json:"popularity"`
		DurationMs  int `json:"duration_ms"`
	} `json:"item"`
	IsPlaying bool `json:"is_playing"`
}

func (c *Spotify) getTrackByPlayerState(ctx context.Context) (*GetTrackResponse, error) {
	data := spotifyPlayerStateResponse{}
	stateResp, err := req.R().
		SetBearerAuthToken(c.channelIntegration.AccessToken).
		SetSuccessResult(&data).
		Get("https://api.spotify.com/v1/me/player")
	if stateResp.StatusCode == 401 && !c.isRetry {
		c.isRetry = true
		c.refreshToken(ctx)
		return c.getTrackByPlayerState(ctx)
	}
	if err != nil {
		return nil, err
	}

	if !stateResp.IsSuccessState() {
		return nil, fmt.Errorf("cannot get player state: %s", stateResp.String())
	}

	var playlistResponse *GetTrackResponsePlaylist
	if data.Context.Type == "playlist" {
		playlistResponse = &GetTrackResponsePlaylist{
			Meta:        nil,
			Uri:         data.Context.URI,
			ExternalUrl: data.Context.ExternalUrls.Spotify,
		}

		splittedUri := strings.Split(data.Context.URI, ":")
		if len(splittedUri) > 2 {
			playlist, err := c.getPlaylist(ctx, splittedUri[2])
			if err != nil {
				return nil, err
			}

			if playlist != nil && playlist.Meta != nil {
				playlistResponse.Meta = playlist.Meta
			}
		}
	}

	artists := make([]string, len(data.Item.Artists))
	for i, artist := range data.Item.Artists {
		artists[i] = artist.Name
	}
	var imageUrl string
	if len(data.Item.Album.Images) > 0 {
		imageUrl = data.Item.Album.Images[0].URL
	}

	return &GetTrackResponse{
		Playlist:  playlistResponse,
		Title:     data.Item.Name,
		Artist:    strings.Join(artists, ", "),
		Image:     imageUrl,
		IsPlaying: data.IsPlaying,
	}, nil
}

type spotifyPlaylistResponse struct {
	Owner struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href        string `json:"href"`
		Id          string `json:"id"`
		Type        string `json:"type"`
		Uri         string `json:"uri"`
		DisplayName string `json:"display_name"`
	} `json:"owner"`
	Followers struct {
		Href  interface{} `json:"href"`
		Total int         `json:"total"`
	} `json:"followers"`
	Description  string `json:"description"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href   string `json:"href"`
	Id     string `json:"id"`
	Name   string `json:"name"`
	Uri    string `json:"uri"`
	Images []struct {
		Height interface{} `json:"height"`
		Width  interface{} `json:"width"`
		Url    string      `json:"url"`
	} `json:"images"`
	Public bool `json:"public"`
}

func (c *Spotify) getPlaylist(ctx context.Context, id string) (*GetTrackResponsePlaylist, error) {
	data := spotifyPlaylistResponse{}
	resp, err := req.R().
		SetBearerAuthToken(c.channelIntegration.AccessToken).
		SetSuccessResult(&data).
		SetQueryParam(
			"fields",
			"description,uri,external_urls,followers,href,id,images,name,owner,public",
		).
		Get(fmt.Sprintf("https://api.spotify.com/v1/playlists/%s", id))
	if resp.StatusCode == 401 && !c.isRetry {
		c.isRetry = true
		c.refreshToken(ctx)
		return c.getPlaylist(ctx, id)
	}
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		if resp.StatusCode == 404 {
			// playlist not found, probably radio by track or something
			return nil, nil
		}

		return nil, fmt.Errorf("cannot get profile: %s", resp.String())
	}

	meta := GetTrackResponsePlaylistMeta{
		Name:      data.Name,
		Images:    []string{},
		Followers: data.Followers.Total,
	}

	for _, image := range data.Images {
		meta.Images = append(meta.Images, image.Url)
	}

	return &GetTrackResponsePlaylist{
		Meta:        &meta,
		Uri:         data.Uri,
		ExternalUrl: data.ExternalUrls.Spotify,
	}, nil
}
