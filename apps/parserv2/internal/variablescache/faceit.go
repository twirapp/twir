package variablescache

import (
	"encoding/json"
	"errors"
	"net/http"
	model "tsuwari/parser/internal/models"

	"github.com/samber/lo"
)

type FaceitGame struct {
	Lvl int `json:"skill_level"`
	Elo int `json:"faceit_elo"`
}

type FaceitResponse struct {
	Games map[string]*FaceitGame `json:"games"`
}

type FaceitDbData struct {
	Game     *string `json:"game"`
	Username string  `json:"username"`
}

func (c *VariablesCacheService) GetFaceitData() (*FaceitGame, error) {
	c.locks.faceitIntegration.Lock()
	defer c.locks.faceitIntegration.Unlock()

	if c.cache.FaceitData != nil {
		return c.cache.FaceitData, nil
	}

	integrations := c.GetEnabledIntegrations()

	if integrations == nil {
		return nil, errors.New("integrations not enabled")
	}

	integration, ok := lo.Find(*integrations, func(i model.ChannelInegrationWithRelation) bool {
		return i.Integration.Service == "FACEIT"
	})

	if !ok {
		return nil, errors.New("faceit integration not enabled")
	}

	dbData := &FaceitDbData{}
	err := json.Unmarshal([]byte(integration.Data.String), &dbData)

	if err != nil {
		return nil, errors.New("failed to read your faceit config. Are you sure you are using integration right?")
	}

	var game string

	if dbData.Game == nil {
		game = "csgo"
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://open.faceit.com/data/v4/players?nickname="+"Satonteu", nil)
	req.Header.Set("Authorization", "Bearer "+integration.Integration.APIKey.String)
	res, err := client.Do(req)

	if err != nil {
		return nil, errors.New("failed to fetch data from faceit")
	}

	data := FaceitResponse{}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, errors.New("failed to fetch data from faceit")
	}

	if data.Games[game] == nil {
		return nil, errors.New("Game " + game + " not found in faceit response.")
	}

	c.cache.FaceitData = data.Games[game]

	return data.Games[game], nil
}
