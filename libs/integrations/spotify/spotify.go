package spotify

import (
	"errors"
	"strings"

	model "github.com/satont/twir/libs/gomodels"

	"github.com/guregu/null"
	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Spotify struct {
	integration *model.ChannelsIntegrations
	isRetry     bool
	db          *gorm.DB
}

func New(integration *model.ChannelsIntegrations, db *gorm.DB) *Spotify {
	if integration == nil || !integration.AccessToken.Valid || !integration.RefreshToken.Valid {
		return nil
	}

	service := Spotify{
		integration: integration,
		db:          db,
		isRetry:     false,
	}

	return &service
}

type SpotifyRefreshResponse struct {
	AccessToken string `json:"access_token"`
}

func (c *Spotify) refreshToken() *error {
	data := SpotifyRefreshResponse{}

	body := make(map[string]string, 2)
	body["grant_type"] = "refresh_token"
	body["refresh_token"] = c.integration.RefreshToken.String

	resp, err := req.R().
		SetFormData(body).
		SetResult(&data).
		SetBasicAuth(
			c.integration.Integration.ClientID.String,
			c.integration.Integration.ClientSecret.String,
		).
		Post("https://accounts.spotify.com/api/token")

	if err != nil || resp.StatusCode != 200 {
		res := errors.New("cannot refresh token")
		return &res
	}

	c.integration.AccessToken = null.StringFrom(data.AccessToken)
	c.db.Where(`"id" = ?`, c.integration.ID).
		Select("*").
		Updates(c.integration).
		Clauses(clause.Returning{})

	return nil
}

type SpotifyArtist struct {
	Name string `json:"name"`
}

type SpotifyImage struct {
	URL string `json:"url"`
}

type SpotifyAlbum struct {
	Images []SpotifyImage `json:"images"`
}

type SpotifyTrack struct {
	Artists []SpotifyArtist `json:"artists"`
	Name    string          `json:"name"`
	Album   SpotifyAlbum    `json:"album"`
}

type SpotifyResponse struct {
	Track     *SpotifyTrack `json:"item"`
	IsPlaying bool          `json:"is_playing"`
}

type GetTrackResponse struct {
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Image     string `json:"image"`
	IsPlaying bool   `json:"isPlaying"`
}

func (c *Spotify) GetTrack() *GetTrackResponse {
	data := SpotifyResponse{}
	req, err := req.R().
		SetBearerAuthToken(c.integration.AccessToken.String).
		SetSuccessResult(&data).
		Get("https://api.spotify.com/v1/me/player/currently-playing")

	if req.StatusCode == 401 && !c.isRetry {
		c.isRetry = true
		c.refreshToken()
		return c.GetTrack()
	}

	if err != nil {
		return nil
	}

	if data.Track == nil {
		return nil
	}

	artistsMap := lo.Map(
		data.Track.Artists, func(artist SpotifyArtist, _ int) string {
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
	}
}

type SpotifyProfile struct {
	Country         string `json:"country"`
	DisplayName     string `json:"display_name"`
	Email           string `json:"email"`
	ExplicitContent struct {
		FilterEnabled bool `json:"filter_enabled"`
		FilterLocked  bool `json:"filter_locked"`
	} `json:"explicit_content"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Href  string `json:"href"`
		Total int    `json:"total"`
	} `json:"followers"`
	Href   string `json:"href"`
	ID     string `json:"id"`
	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	Product string `json:"product"`
	Type    string `json:"type"`
	URI     string `json:"uri"`
}

func (c *Spotify) GetProfile() (*SpotifyProfile, error) {
	data := SpotifyProfile{}
	req, err := req.R().
		SetBearerAuthToken(c.integration.AccessToken.String).
		SetSuccessResult(&data).
		Get("https://api.spotify.com/v1/me")

	if req.StatusCode == 401 && !c.isRetry {
		c.isRetry = true
		c.refreshToken()
		return c.GetProfile()
	}

	if err != nil {
		return nil, err
	}

	if !req.IsSuccessState() {
		return nil, errors.New("cannot get profile")
	}

	return &data, nil
}
