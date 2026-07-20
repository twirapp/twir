package gsi

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

type Player struct {
	SteamID   string         `json:"steamid"`
	AccountID int64          `json:"account_id"`
	Name      string         `json:"name"`
	Activity  PlayerActivity `json:"activity"`
	Kills     int            `json:"kills"`
	Deaths    int            `json:"deaths"`
	Assists   int            `json:"assists"`
	TeamName  string         `json:"team_name"`
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
	GameTime   int    `json:"game_time"`
}

type Payload struct {
	Provider Provider `json:"provider"`
	Auth     Auth     `json:"auth"`
	Map      *Map     `json:"map"`
	Player   *Player  `json:"player"`
	Hero     *Hero    `json:"hero"`
	Events   []Event  `json:"events"`
}
