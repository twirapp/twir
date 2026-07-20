package api

const DotaStateUpdateSubject = "api.dota.state_update"

type DotaStateUpdateMessage struct {
	ChannelID      string  `json:"channelId"`
	InGame         bool    `json:"inGame"`
	Mmr            int     `json:"mmr"`
	SessionWins    int     `json:"sessionWins"`
	SessionLosses  int     `json:"sessionLosses"`
	WinProbability float64 `json:"winProbability"`
	HeroName       string  `json:"heroName"`
	MatchID        int64   `json:"matchId"`
}
