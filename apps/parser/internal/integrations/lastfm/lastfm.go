package lastfm

import (
	"fmt"

	model "github.com/satont/tsuwari/libs/gomodels"

	req "github.com/imroc/req/v3"
)

type LastFm struct {
	integration *model.ChannelsIntegrations
}

func New(integration *model.ChannelsIntegrations) *LastFm {
	if integration == nil || integration.Data == nil || !integration.Integration.APIKey.Valid {
		return nil
	}

	service := LastFm{
		integration: integration,
	}

	return &service
}

type LastFmTrack struct {
	Artist struct {
		Text string `json:"#text"`
	} `json:"artist"`
	Attr *struct {
		NowPlaying *string `json:"nowplaying"`
	} `json:"@attr"`
	Album *struct {
		Text string `json:"#text"`
	} `json:"album"`
	Name string `json:"name"`
}

type LastFmResponse struct {
	Error   *int
	Message *string

	RecentTracks *struct {
		Track *[]*LastFmTrack `json:"track"`
	} `json:"recenttracks"`
}

func (c *LastFm) GetTrack() *string {
	data := LastFmResponse{}
	var response string

	resp, err := req.R().
		SetQueryParam("method", "user.getrecenttracks").
		SetQueryParam("user", *c.integration.Data.UserName).
		SetQueryParam("api_key", c.integration.Integration.APIKey.String).
		SetQueryParam("format", "json").
		SetQueryParam("limit", "1").
		SetResult(&data).
		SetContentType("application/json").
		Get("http://ws.audioscrobbler.com/2.0")

	if err != nil || !resp.IsSuccess() {
		return nil
	}

	if data.RecentTracks == nil || data.RecentTracks.Track == nil {
		return nil
	}
	tracks := *data.RecentTracks.Track
	track := tracks[0]
	if track == nil || track.Attr == nil || track.Attr.NowPlaying == nil {
		return nil
	}

	response = fmt.Sprintf("%s — %s", track.Artist.Text, track.Name)

	return &response
}
