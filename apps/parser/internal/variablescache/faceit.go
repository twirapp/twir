package variables_cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	model "tsuwari/models"

	"github.com/samber/lo"
)

type FaceitResult struct {
	FaceitUser *FaceitUser
	Matches    *[]FaceitMatch `json:"matches"`
}

type FaceitGame struct {
	Name string
	Lvl  int `json:"skill_level"`
	Elo  int `json:"faceit_elo"`
}

type FaceitUser struct {
	FaceitGame
	PlayerId string
}

type FaceitUserResponse struct {
	PlayerId string                 `json:"player_id"`
	Games    map[string]*FaceitGame `json:"games"`
}

func (c *VariablesCacheService) GetFaceitUserData() (*FaceitUser, error) {
	c.locks.faceitIntegration.Lock()
	defer c.locks.faceitIntegration.Unlock()

	if c.cache.FaceitData != nil && c.cache.FaceitData.FaceitUser != nil {
		return c.cache.FaceitData.FaceitUser, nil
	}

	c.cache.FaceitData = &FaceitResult{}

	integrations := c.GetEnabledIntegrations()

	if integrations == nil {
		return nil, errors.New("no enabled integrations")
	}

	integration, ok := lo.Find(*integrations, func(i model.ChannelInegrationWithRelation) bool {
		return i.Integration.Service == "FACEIT" && i.Enabled
	})

	if !ok {
		return nil, errors.New("faceit integration not enabled")
	}

	var game string = *integration.Data.Game

	if integration.Data.Game == nil {
		game = "csgo"
	}

	client := &http.Client{}
	req, _ := http.NewRequest(
		"GET",
		"https://open.faceit.com/data/v4/players?nickname="+*integration.Data.UserName,
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+integration.Integration.APIKey.String)
	res, err := client.Do(req)

	if req.Response != nil && req.Response.StatusCode == 404 {
		return nil, errors.New(
			"user not found on faceit. Please make sure you typed correct nickname",
		)
	}

	if err != nil {
		return nil, err
	}

	data := FaceitUserResponse{}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, errors.New("internal error happend on parsing user profile.")
	}

	if data.Games[game] == nil {
		return nil, errors.New(game + " game not found in faceit response.")
	}

	data.Games[game].Name = game

	c.cache.FaceitData.FaceitUser = &FaceitUser{
		FaceitGame: *data.Games[game],
		PlayerId:   data.PlayerId,
	}

	return c.cache.FaceitData.FaceitUser, nil
}

type FaceitMatch struct {
	Team         string  `json:"i5"`
	TeamScore    string  `json:"i18"`
	Map          string  `json:"i1"`
	Kd           string  `json:"c2"`
	HsPercentage string  `json:"c4"`
	HsNumber     string  `json:"i13"`
	Kills        string  `json:"i6"`
	Deaths       string  `json:"i8"`
	CreatedAt    int64   `json:"created_at"`
	UpdateAt     int64   `json:"updated_at"`
	Elo          *string `json:"elo"`
	EloDiff      *string
	IsWin        bool

	RawIsWin string `json:"i10"`
}

type FaceitMatchesResponse []FaceitMatch

func (c *VariablesCacheService) GetFaceitLatestMatches() (*[]FaceitMatch, error) {
	c.locks.faceitMatches.Lock()
	defer c.locks.faceitMatches.Unlock()

	if c.cache.FaceitData == nil || c.cache.FaceitData.FaceitUser == nil {
		faceitUser, err := c.GetFaceitUserData()
		if err != nil {
			return nil, err
		}

		c.cache.FaceitData.FaceitUser = faceitUser
	}

	if c.cache.FaceitData.Matches != nil {
		return c.cache.FaceitData.Matches, nil
	}

	client := &http.Client{}

	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"https://api.faceit.com/stats/api/v1/stats/time/users/%s/games/%s?size=30",
			c.cache.FaceitData.FaceitUser.PlayerId,
			c.cache.FaceitData.FaceitUser.FaceitGame.Name,
		),
		nil,
	)
	res, _ := client.Do(req)

	reqResult := FaceitMatchesResponse{}
	json.NewDecoder(res.Body).Decode(&reqResult)

	matches := []FaceitMatch{}
	stream := c.GetChannelStream()
	if stream == nil {
		return &matches, nil
	}
	startedDate := stream.StartedAt.UnixMilli()

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

		prevMatch := &reqResult[i+1]
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

	return &matches, nil
}

func (c *VariablesCacheService) GetFaceitTodayEloDiff(matches *[]FaceitMatch) int {
	if matches == nil {
		return 0
	}

	sum := lo.Reduce(*matches, func(agg int, item FaceitMatch, _ int) int {
		if item.EloDiff == nil {
			return agg
		}
		v, err := strconv.Atoi(*item.EloDiff)
		if err != nil {
			return agg
		}
		return agg + v
	}, 0)

	return sum
}
