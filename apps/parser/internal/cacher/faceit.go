package cacher

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
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

	reqResult := faceitMatchesResponse{}

	_, err = req.C().EnableForceHTTP1().R().
		SetContext(ctx).
		SetSuccessResult(&reqResult).
		Get(
			fmt.Sprintf(
				"https://api.faceit.com/stats/api/v1/stats/time/users/%s/games/%s?size=30",
				c.cache.faceitData.FaceitUser.PlayerId,
				c.cache.faceitData.FaceitUser.FaceitGame.Name,
			),
		)

	if err != nil {
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

	data := &types.FaceitUserResponse{}
	resp, err := req.C().EnableForceHTTP1().R().
		SetContext(ctx).
		SetBearerAuthToken(integration.Integration.APIKey.String).
		SetSuccessResult(data).
		Get("https://open.faceit.com/data/v4/players/" + *integration.Data.UserId)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, errors.New(
			"user not found on faceit. Please make sure you typed correct nickname",
		)
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

	var stateData faceitStateResponse
	resp, err := req.C().EnableForceHTTP1().R().
		SetContext(ctx).
		SetSuccessResult(&stateData).
		Get("https://api.faceit.com/match/v1/matches/groupByState?userId=" + faceitUser.PlayerId)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, errors.New("cannot get data from faceit")
	}

	// match not going
	if stateData.Payload.Ongoing.ID == "" {
		return nil, nil
	}

	var matchData faceitMatchResponse
	matchResp, err := req.C().EnableForceHTTP1().R().
		SetContext(ctx).
		SetSuccessResult(&matchData).
		Get("https://api.faceit.com/match/v2/match/" + stateData.Payload.Ongoing.ID)
	if err != nil {
		return nil, err
	}
	if !matchResp.IsSuccessState() {
		return nil, errors.New("cannot get data from faceit")
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
