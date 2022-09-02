package variablescache

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	model "tsuwari/parser/internal/models"

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

type FaceitDbData struct {
	Game     *string `json:"game"`
	Username string  `json:"username"`
}

func (c *VariablesCacheService) GetFaceitUserData() *FaceitUser {
	c.locks.faceitIntegration.Lock()
	defer c.locks.faceitIntegration.Unlock()

	if c.cache.FaceitData != nil && c.cache.FaceitData.FaceitUser != nil {
		return c.cache.FaceitData.FaceitUser
	}

	c.cache.FaceitData = &FaceitResult{}

	integrations := c.GetEnabledIntegrations()

	if integrations == nil {
		return nil
	}

	integration, ok := lo.Find(*integrations, func(i model.ChannelInegrationWithRelation) bool {
		return i.Integration.Service == "FACEIT"
	})

	if !ok {
		return nil
	}

	dbData := &FaceitDbData{}
	err := json.Unmarshal([]byte(integration.Data.String), &dbData)

	if err != nil {
		return nil
	}

	var game string

	if dbData.Game == nil {
		game = "csgo"
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://open.faceit.com/data/v4/players?nickname="+dbData.Username, nil)
	req.Header.Set("Authorization", "Bearer "+integration.Integration.APIKey.String)
	res, err := client.Do(req)

	if err != nil {
		return nil
	}

	data := FaceitUserResponse{}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil
	}

	if data.Games[game] == nil {
		return nil
	}

	data.Games[game].Name = game

	c.cache.FaceitData.FaceitUser = &FaceitUser{
		FaceitGame: *data.Games[game],
		PlayerId:   data.PlayerId,
	}

	return c.cache.FaceitData.FaceitUser
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

func (c *VariablesCacheService) GetFaceitLatestMatches() *[]FaceitMatch {
	c.locks.faceitMatches.Lock()
	defer c.locks.faceitMatches.Unlock()

	if c.cache.FaceitData == nil || c.cache.FaceitData.FaceitUser == nil {
		c.cache.FaceitData.FaceitUser = c.GetFaceitUserData()
	}

	if c.cache.FaceitData.Matches != nil {
		return c.cache.FaceitData.Matches
	}

	client := &http.Client{}
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.faceit.com/stats/api/v1/stats/time/users/%s/games/%s?size=30", c.cache.FaceitData.FaceitUser.PlayerId, c.cache.FaceitData.FaceitUser.FaceitGame.Name),
		nil,
	)
	res, _ := client.Do(req)

	reqResult := FaceitMatchesResponse{}
	json.NewDecoder(res.Body).Decode(&reqResult)

	matches := []FaceitMatch{}
	stream := c.GetChannelStream()
	if stream == nil {
		return &matches
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

	return &matches
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
