package lastfm

import (
	"encoding/json"
	"fmt"
	model "tsuwari/parser/internal/models"

	req "github.com/imroc/req/v3"
)

type DbData struct {
	UserName *string `json:"username"`
}

type LastFm struct {
	DbData
	integration *model.ChannelInegrationWithRelation
}

func New(integration *model.ChannelInegrationWithRelation) *LastFm {
	if integration == nil || !integration.Data.Valid || !integration.Integration.APIKey.Valid {
		return nil
	}

	dbData := DbData{}
	err := json.Unmarshal([]byte(integration.Data.String), &dbData)
	if err != nil {
		return nil
	}

	service := LastFm{
		DbData:      dbData,
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

func (c *LastFm) GetRecentTrack() *string {
	data := LastFmResponse{}
	var response string

	resp, err := req.R().
		SetQueryParam("method", "user.getrecenttracks").
		SetQueryParam("user", *c.DbData.UserName).
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

/* func (c *LastFm) GetRecentTrack() *string {
	data := LastFmResponse{}
	var response string

	rBuilder := requests.
		URL("http://ws.audioscrobbler.com/2.0").
		ContentType("application/json").
		ToJSON(&data)

	rBuilder.Param("method", "user.getrecenttracks")
	rBuilder.Param("user", *c.DbData.UserName)
	rBuilder.Param("api_key", c.integration.Integration.APIKey.String)
	rBuilder.Param("format", "json")
	rBuilder.Param("limit", "1")

	err := rBuilder.Fetch(context.Background())

	if err != nil {
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
*/
