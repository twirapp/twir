package types

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

type FaceitEstimateGainLose struct {
	Gain int
	Lose int
}

type FaceitResult struct {
	FaceitUser   *FaceitUser
	Matches      []*FaceitMatch `json:"matches"`
	EstimateGain int
	EstimateLose int
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
