package gsi

import (
	"encoding/json"
	"strconv"
)

type GameState string

const (
	GameStateHeroSelection GameState = "DOTA_GAMERULES_STATE_HERO_SELECTION"
	GameStateStrategyTime  GameState = "DOTA_GAMERULES_STATE_STRATEGY_TIME"
	GameStatePreGame       GameState = "DOTA_GAMERULES_STATE_PRE_GAME"
	GameStateInProgress    GameState = "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS"
	GameStatePostGame      GameState = "DOTA_GAMERULES_STATE_POST_GAME"
)

type WinTeam string

const (
	WinTeamRadiant WinTeam = "radiant"
	WinTeamDire    WinTeam = "dire"
	WinTeamNone    WinTeam = "none"
)

type PlayerActivity string

const (
	PlayerActivityPlaying PlayerActivity = "playing"
)

type Provider struct {
	Name      string `json:"name"`
	AppID     int    `json:"appid"`
	Version   int    `json:"version"`
	Timestamp int64  `json:"timestamp"`
}

type Auth struct {
	Token string `json:"token"`
}

type Map struct {
	Name         string    `json:"name"`
	MatchID      int64     `json:"matchid"`
	GameTime     int       `json:"game_time"`
	ClockTime    int       `json:"clock_time"`
	GameState    GameState `json:"game_state"`
	Paused       bool      `json:"paused"`
	WinTeam      WinTeam   `json:"win_team"`
	RadiantScore int       `json:"radiant_score"`
	DireScore    int       `json:"dire_score"`
}

// steamID64Base is the SteamID64 individual account base; the 32-bit Dota
// account ID is steamID64 - steamID64Base.
const steamID64Base int64 = 76561197960265728

type Player struct {
	SteamID string `json:"steamid"`
	// Real GSI payloads carry the account id as "accountid" (a quoted decimal
	// string in observed captures); "account_id" does not exist in the wire
	// format. json.Number accepts both quoted and unquoted forms.
	AccountIDRaw json.Number    `json:"accountid"`
	Name         string         `json:"name"`
	Activity     PlayerActivity `json:"activity"`
	Kills        int            `json:"kills"`
	Deaths       int            `json:"deaths"`
	Assists      int            `json:"assists"`
	TeamName     string         `json:"team_name"`
}

// DotaAccountID resolves the 32-bit Dota account ID, preferring the
// always-present SteamID64 and falling back to the raw "accountid" field.
// Returns 0 when neither source is usable.
func (p Player) DotaAccountID() int64 {
	if steamID, err := strconv.ParseInt(p.SteamID, 10, 64); err == nil && steamID >= steamID64Base {
		return steamID - steamID64Base
	}
	if accountID, err := p.AccountIDRaw.Int64(); err == nil && accountID > 0 {
		return accountID
	}

	return 0
}

type Hero struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
}

type Event struct {
	EventType  string `json:"event_type"`
	KillerTeam string `json:"killer_team,omitempty"`
	PlayerID   *int   `json:"player_id,omitempty"`
	Player     *int   `json:"player,omitempty"`
	GameTime   int    `json:"game_time"`
}

func (e Event) AegisPlayerID() *int {
	if e.PlayerID != nil {
		return e.PlayerID
	}

	return e.Player
}

type Payload struct {
	Provider Provider `json:"provider"`
	Auth     Auth     `json:"auth"`
	Map      *Map     `json:"map"`
	Player   *Player  `json:"player"`
	Hero     *Hero    `json:"hero"`
	Events   []Event  `json:"events"`
}
