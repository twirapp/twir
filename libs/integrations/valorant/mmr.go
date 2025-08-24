package valorant

import (
	"context"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
)

type MmrResponse struct {
	Status int              `json:"status"`
	Data   *MmrResponseData `json:"data"`
}

type MmrResponseData struct {
	Account struct {
		Puuid string `json:"puuid"`
		Name  string `json:"name"`
		Tag   string `json:"tag"`
	} `json:"account"`
	Peak struct {
		Season struct {
			Id    string `json:"id"`
			Short string `json:"short"`
		} `json:"season"`
		RankingSchema string `json:"ranking_schema"`
		Tier          struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"tier"`
	} `json:"peak"`
	Current struct {
		Tier struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"tier"`
		Rr                   int `json:"rr"`
		LastChange           int `json:"last_change"`
		Elo                  int `json:"elo"`
		GamesNeededForRating int `json:"games_needed_for_rating"`
		LeaderboardPlacement struct {
			Rank      int       `json:"rank"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"leaderboard_placement"`
	} `json:"current"`
	Seasonal []struct {
		Season struct {
			Id    string `json:"id"`
			Short string `json:"short"`
		} `json:"season"`
		Wins    int `json:"wins"`
		Games   int `json:"games"`
		EndTier struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"end_tier"`
		RankingSchema        string `json:"ranking_schema"`
		LeaderboardPlacement struct {
			Rank      int       `json:"rank"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"leaderboard_placement"`
		ActWins []struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"act_wins"`
	} `json:"seasonal"`
}

func (c *HenrikValorantApiClient) GetValorantProfileMmr(
	ctx context.Context,
	platform,
	region,
	puuid string,
) (
	*MmrResponse,
	error,
) {
	apiUrl := fmt.Sprintf(
		"https://api.henrikdev.xyz/valorant/v3/by-puuid/mmr/%s/%s/%s",
		region,
		platform,
		puuid,
	)

	data := &MmrResponse{}

	response, err := req.R().
		SetContext(ctx).
		SetHeader("Authorization", c.apiKey).
		SetSuccessResult(data).Get(apiUrl)
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf(
			"cannot get valorant profile for puuid %s in region %s: %s",
			puuid,
			region,
			response.String(),
		)
	}

	return data, nil
}
