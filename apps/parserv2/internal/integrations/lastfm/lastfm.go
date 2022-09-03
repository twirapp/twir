package lastfm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	model "tsuwari/parser/internal/models"
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
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://ws.audioscrobbler.com/2.0", nil)
	q := req.URL.Query()
	q.Add("method", "user.getrecenttracks")
	q.Add("user", *c.DbData.UserName)
	q.Add("api_key", c.integration.Integration.APIKey.String)
	q.Add("format", "json")
	q.Add("limit", "1")
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)

	if err != nil {
		return nil
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println("body", bodyString)

	var response string
	data := LastFmResponse{}
	err = json.Unmarshal(bodyBytes, &data)
	if err == nil {
		if data.RecentTracks == nil || data.RecentTracks.Track == nil {
			return nil
		}
		tracks := *data.RecentTracks.Track
		track := tracks[0]
		if track == nil || track.Attr == nil || track.Attr.NowPlaying == nil {
			return nil
		}

		response = fmt.Sprintf("%s — %s", track.Artist.Text, track.Name)
	}

	return &response
}
