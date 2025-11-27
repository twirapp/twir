package cacher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
)

type faceitMatchesResponse []*types.FaceitMatch

// GetFaceitLatestMatches implements types.VariablesCacher
func (c *cacher) GetFaceitLatestMatches(ctx context.Context) ([]*types.FaceitMatch, error) {
	c.locks.faceitMatches.Lock()
	defer c.locks.faceitMatches.Unlock()

	_, err := c.GetFaceitUserData(ctx)
	if err != nil {
		return nil, err
	}

	if c.cache.faceitData.Matches != nil {
		return c.cache.faceitData.Matches, nil
	}

	reqUrl := fmt.Sprintf(
		"https://api.faceit.com/stats/api/v1/stats/time/users/%s/games/%s?size=30",
		c.cache.faceitData.FaceitUser.PlayerId,
		c.cache.faceitData.FaceitUser.FaceitGame.Name,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil, err
	}

	var reqResult faceitMatchesResponse
	if err := json.Unmarshal(body, &reqResult); err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil, err
	}

	var matches []*types.FaceitMatch
	// TODO: use stream time
	// stream := c.GetChannelStream(ctx)
	// if stream == nil {
	// 	return matches, nil
	// }
	// startedDate := stream.StartedAt.UnixMilli()
	startedDate := time.Now().UTC().Truncate(24*time.Hour).UnixMilli() + 1000

	for i, match := range reqResult {
		matchCreatedAt := time.UnixMilli(match.UpdateAt).UnixMilli()

		if matchCreatedAt < startedDate {
			continue
		}

		if i+1 > len(reqResult)-1 {
			break
		}

		val := false
		if match.RawIsWin == "1" {
			val = true
		}
		match.IsWin = val

		if i+1 >= len(reqResult)-1 {
			break
		}

		prevMatch := reqResult[i+1]
		if prevMatch == nil || prevMatch.Elo == nil || match.Elo == nil {
			continue
		}

		prevElo, pErr := strconv.Atoi(*prevMatch.Elo)
		currElo, cErr := strconv.Atoi(*match.Elo)

		if pErr != nil || cErr != nil {
			continue
		}

		var eloDiff int
		if *prevMatch.Elo > *match.Elo {
			eloDiff = -(prevElo - currElo)
		} else {
			eloDiff = currElo - prevElo
		}

		newMatchEloDiff := strconv.Itoa(eloDiff)
		match.EloDiff = &newMatchEloDiff
		matches = append(matches, match)
	}

	c.cache.faceitData.Matches = matches

	return matches, nil
}

// GetFaceitTodayEloDiff implements types.VariablesCacher
func (c *cacher) GetFaceitTodayEloDiff(_ context.Context, matches []*types.FaceitMatch) int {
	if matches == nil {
		return 0
	}

	sum := lo.Reduce(
		matches, func(agg int, item *types.FaceitMatch, _ int) int {
			if item.EloDiff == nil {
				return agg
			}
			v, err := strconv.Atoi(*item.EloDiff)
			if err != nil {
				return agg
			}
			return agg + v
		}, 0,
	)

	return sum
}

// GetFaceitUserData implements types.VariablesCacher
func (c *cacher) GetFaceitUserData(ctx context.Context) (*types.FaceitUser, error) {
	c.locks.faceitUserData.Lock()
	defer c.locks.faceitUserData.Unlock()

	if c.cache.faceitData != nil && c.cache.faceitData.FaceitUser != nil {
		return c.cache.faceitData.FaceitUser, nil
	}

	c.cache.faceitData = &types.FaceitResult{}

	integrations := c.GetEnabledChannelIntegrations(ctx)

	if integrations == nil {
		return nil, errors.New("no enabled integrations")
	}

	integration, ok := lo.Find(
		integrations, func(i *model.ChannelsIntegrations) bool {
			return i.Integration.Service == "FACEIT" && i.Enabled
		},
	)

	if !ok {
		return nil, errors.New("faceit integration not enabled")
	}

	var game = *integration.Data.Game

	if integration.Data.Game == nil {
		game = "cs2"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://open.faceit.com/data/v4/players/"+*integration.Data.UserId, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+integration.Integration.APIKey.String)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, errors.New(
			"user not found on faceit. Please make sure you typed correct nickname",
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data types.FaceitUserResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if data.Games[game] == nil {
		return nil, errors.New(game + " game not found in faceit response.")
	}

	data.Games[game].Name = game

	c.cache.faceitData.FaceitUser = &types.FaceitUser{
		FaceitGame: *data.Games[game],
		PlayerId:   data.PlayerId,
	}

	return c.cache.faceitData.FaceitUser, nil
}

type faceitStateResponse struct {
	Payload struct {
		Ongoing struct {
			ID string `json:"id"`
		} `json:"ONGOING"`
	} `json:"payload"`
}

type faceitMatchResponseFaction struct {
	Roster []struct {
		ID string `json:"id"`
	} `json:"roster"`
	Stats struct {
		WinProbability float64 `json:"winProbability"`
	} `json:"stats"`
}

type faceitMatchResponse struct {
	Payload struct {
		Teams struct {
			Faction1 faceitMatchResponseFaction `json:"faction1"`
			Faction2 faceitMatchResponseFaction `json:"faction2"`
		} `json:"teams"`
	} `json:"payload"`
}

func (c *cacher) ComputeFaceitGainLoseEstimate(ctx context.Context) (
	*types.FaceitEstimateGainLose,
	error,
) {
	faceitUser, err := c.GetFaceitUserData(ctx)
	if err != nil {
		return nil, err
	}

	if c.cache.faceitData.EstimateLose != 0 || c.cache.faceitData.EstimateGain != 0 {
		return &types.FaceitEstimateGainLose{
			Lose: c.cache.faceitData.EstimateLose,
			Gain: c.cache.faceitData.EstimateGain,
		}, nil
	}

	stateReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.faceit.com/match/v1/matches/groupByState?userId="+faceitUser.PlayerId, nil)
	if err != nil {
		return nil, err
	}

	stateResp, err := http.DefaultClient.Do(stateReq)
	if err != nil {
		return nil, err
	}
	defer stateResp.Body.Close()

	if stateResp.StatusCode < 200 || stateResp.StatusCode >= 300 {
		return nil, errors.New("cannot get data from faceit")
	}

	stateBody, err := io.ReadAll(stateResp.Body)
	if err != nil {
		return nil, err
	}

	var stateData faceitStateResponse
	if err := json.Unmarshal(stateBody, &stateData); err != nil {
		return nil, err
	}

	// match not going
	if stateData.Payload.Ongoing.ID == "" {
		return nil, nil
	}

	matchReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.faceit.com/match/v2/match/"+stateData.Payload.Ongoing.ID, nil)
	if err != nil {
		return nil, err
	}

	matchResp, err := http.DefaultClient.Do(matchReq)
	if err != nil {
		return nil, err
	}
	defer matchResp.Body.Close()

	if matchResp.StatusCode < 200 || matchResp.StatusCode >= 300 {
		return nil, errors.New("cannot get data from faceit")
	}

	matchBody, err := io.ReadAll(matchResp.Body)
	if err != nil {
		return nil, err
	}

	var matchData faceitMatchResponse
	if err := json.Unmarshal(matchBody, &matchData); err != nil {
		return nil, err
	}

	// by default set user team to faction 1, if user found in faction2, then we reassign team
	playerTeam := matchData.Payload.Teams.Faction1
	for _, p := range matchData.Payload.Teams.Faction2.Roster {
		if p.ID == faceitUser.PlayerId {
			playerTeam = matchData.Payload.Teams.Faction2
			break
		}
	}

	winProbability := playerTeam.Stats.WinProbability
	gain := math.Round(50 - winProbability*50)

	c.locks.faceitUserData.Lock()
	defer c.locks.faceitUserData.Unlock()

	c.cache.faceitData.EstimateGain = int(gain)
	c.cache.faceitData.EstimateLose = int(50 - gain)

	return &types.FaceitEstimateGainLose{
		Lose: c.cache.faceitData.EstimateLose,
		Gain: c.cache.faceitData.EstimateGain,
	}, nil
}
