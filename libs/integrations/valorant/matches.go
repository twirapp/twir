package valorant

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type StoredMatchesResponse struct {
	Status  int    `json:"status"`
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	Results struct {
		Total    int `json:"total"`
		Returned int `json:"returned"`
		Before   int `json:"before"`
		After    int `json:"after"`
	} `json:"results"`
	Data []StoredMatchesResponseMatch `json:"data"`
}

type StoredMatchesResponseMatch struct {
	Meta struct {
		Id  string `json:"id"`
		Map struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"map"`
		Version   string                   `json:"version"`
		Mode      StoreMatchesResponseMode `json:"mode"`
		StartedAt time.Time                `json:"started_at"`
		Season    struct {
			Id    string `json:"id"`
			Short string `json:"short"`
		} `json:"season"`
		Region  string `json:"region"`
		Cluster string `json:"cluster"`
	} `json:"meta"`
	Stats StoredMatchesResponseMatchStats `json:"stats"`
	Teams struct {
		Red  int `json:"red"`
		Blue int `json:"blue"`
	} `json:"teams"`
}

type StoredMatchesResponseMatchStats struct {
	Puuid     string `json:"puuid"`
	Team      string `json:"team"`
	Level     int    `json:"level"`
	Character struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"character"`
	Tier    int `json:"tier"`
	Score   int `json:"score"`
	Kills   int `json:"kills"`
	Deaths  int `json:"deaths"`
	Assists int `json:"assists"`
	Shots   struct {
		Head int `json:"head"`
		Body int `json:"body"`
		Leg  int `json:"leg"`
	} `json:"shots"`
	Damage struct {
		Dealt    int `json:"dealt"`
		Received int `json:"received"`
	} `json:"damage"`
}

type StoreMatchesResponseMode string

const (
	StoreMatchesResponseModeCompetitive StoreMatchesResponseMode = "Competitive"
	StoreMatchesResponseModeUnrated     StoreMatchesResponseMode = "Unrated"
)

func (c *HenrikValorantApiClient) GetProfileStoredMatches(
	ctx context.Context,
	region,
	puuid string,
) (*StoredMatchesResponse, error) {
	apiUrl := fmt.Sprintf(
		"https://api.henrikdev.xyz/valorant/v1/by-puuid/stored-matches/%s/%s?size=100",
		region,
		puuid,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"cannot get valorant stored matches for puuid %s in region %s: %s",
			puuid,
			region,
			string(body),
		)
	}

	var data *StoredMatchesResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return data, nil
}
