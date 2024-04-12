package types

import "github.com/NovikovRoman/pubg"

type PubgData struct {
	LifetimeStats pubg.LifetimeStatsPlayer
	SeasonStats   pubg.SeasonStatsPlayer
	Mastery       pubg.SurvivalMastery
}

type PubgSeason struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes struct {
		IsCurrentSeason bool `json:"isCurrentSeason"`
		IsOffseason     bool `json:"isOffseason"`
	} `json:"attributes"`
}
