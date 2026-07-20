package dota

const (
	GetDataSubject      = "dota.get_data"
	MatchStartedSubject = "dota.match_started"
	MatchEndedSubject   = "dota.match_ended"
	RoshanKilledSubject = "dota.roshan_killed"
	AegisPickupSubject  = "dota.aegis_pickup"
)

type GetDataRequest struct {
	ChannelID      string
	TwitchUserID   string
	SteamAccountID string
}

type LastGameInfo struct {
	HeroName  string
	Kills     int
	Deaths    int
	Assists   int
	Win       bool
	DurationS int
}

type GetDataResponse struct {
	Enabled        bool
	Linked         bool
	InGame         bool
	Mmr            int
	SessionWins    int
	SessionLosses  int
	HeroName       string
	MatchID        int64
	TeamIsRadiant  bool
	RadiantScore   int
	DireScore      int
	GameTime       int
	WinProbability float64
	NotablePlayers []string
	LastGame       *LastGameInfo
}

type MatchStartedMessage struct {
	ChannelID      string
	TwitchUserID   string
	SteamAccountID string
	HeroName       string
}

type MatchEndedMessage struct {
	ChannelID      string
	TwitchUserID   string
	SteamAccountID string
	Win            bool
	HeroName       string
	Mmr            int
	SessionWins    int
	SessionLosses  int
}

type RoshanKilledMessage struct {
	ChannelID    string
	TwitchUserID string
	Team         string
	GameTime     int
}

type AegisPickupMessage struct {
	ChannelID    string
	TwitchUserID string
	PlayerName   string
	GameTime     int
}
