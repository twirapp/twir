package types

type RocketLeagueUserStats struct {
	Rankings []Ranking `json:"ranking"`
	Username string    `json:"username"`
}

type Ranking struct {
	Playlist string   `json:"playlist"`
	Rating   int32    `json:"rating"`
	Rank     string   `json:"rank"`
	Matches  *Matches `json:"matches"`
}

type Matches struct {
	Total  int32  `json:"total"`
	Streak string `json:"streak"`
}
