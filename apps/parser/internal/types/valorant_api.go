package types

import (
	"time"
)

type ValorantProfileImages struct {
	Small        string `json:"small"`
	Large        string `json:"large"`
	TriangleDown string `json:"triangle_down"`
	TriangleUp   string `json:"triangle_up"`
}

type ValorantMMR struct {
	Status int `json:"status"`
	Data   struct {
		Account struct {
			Puuid string `json:"puuid"`
			Name  string `json:"name"`
			Tag   string `json:"tag"`
		} `json:"account"`
		Peak struct {
			Season struct {
				Id    string `json:"id"`
				Short string `json:"short"`
			} `json:"season"`
			RankingSchema string `json:"ranking_schema"`
			Tier          struct {
				Id   int    `json:"id"`
				Name string `json:"name"`
			} `json:"tier"`
		} `json:"peak"`
		Current struct {
			Tier struct {
				Id   int    `json:"id"`
				Name string `json:"name"`
			} `json:"tier"`
			Rr                   int `json:"rr"`
			LastChange           int `json:"last_change"`
			Elo                  int `json:"elo"`
			GamesNeededForRating int `json:"games_needed_for_rating"`
			LeaderboardPlacement struct {
				Rank      int       `json:"rank"`
				UpdatedAt time.Time `json:"updated_at"`
			} `json:"leaderboard_placement"`
		} `json:"current"`
		Seasonal []struct {
			Season struct {
				Id    string `json:"id"`
				Short string `json:"short"`
			} `json:"season"`
			Wins    int `json:"wins"`
			Games   int `json:"games"`
			EndTier struct {
				Id   int    `json:"id"`
				Name string `json:"name"`
			} `json:"end_tier"`
			RankingSchema        string `json:"ranking_schema"`
			LeaderboardPlacement struct {
				Rank      int       `json:"rank"`
				UpdatedAt time.Time `json:"updated_at"`
			} `json:"leaderboard_placement"`
			ActWins []struct {
				Id   int    `json:"id"`
				Name string `json:"name"`
			} `json:"act_wins"`
		} `json:"seasonal"`
	} `json:"data"`
}

type ValorantMatchPlayer struct {
	Puuid              string `json:"puuid"`
	Name               string `json:"name"`
	Tag                string `json:"tag"`
	Team               string `json:"team"`
	Character          string `json:"character"`
	CurrentTierPatched string `json:"currenttier_patched"`
	Stats              struct {
		Score     int `json:"score"`
		Kills     int `json:"kills"`
		Deaths    int `json:"deaths"`
		Assists   int `json:"assists"`
		Bodyshots int `json:"bodyshots"`
		Headshots int `json:"headshots"`
		Legshots  int `json:"legshots"`
	} `json:"stats"`
	Behavior struct {
		AfkRounds    float64 `json:"afk_rounds"`
		FriendlyFire struct {
			Incoming float64 `json:"incoming"`
			Outgoing float64 `json:"outgoing"`
		} `json:"friendly_fire"`
		RoundsInSpawn float64 `json:"rounds_in_spawn"`
	} `json:"behavior"`
	Economy struct {
		Spent struct {
			Overall float64 `json:"overall"`
			Average float64 `json:"average"`
		} `json:"spent"`
		LoadoutValue struct {
			Overall float64 `json:"overall"`
			Average float64 `json:"average"`
		} `json:"loadout_value"`
	} `json:"economy"`
	Level          float64 `json:"level"`
	CurrentTier    float64 `json:"currenttier"`
	DamageMade     float64 `json:"damage_made"`
	DamageReceived float64 `json:"damage_received"`
}

type ValorantMatchPlayers struct {
	AllPlayers []ValorantMatchPlayer `json:"all_players"`
}

type ValorantMatchesResponse struct {
	Data []ValorantMatch `json:"data"`
}

type ValorantMatch struct {
	Teams map[string]struct {
		HasWon     bool `json:"has_won"`
		RoundsWon  int  `json:"rounds_won"`
		RoundsLost int  `json:"rounds_lost"`
	}
	Players  ValorantMatchPlayers `json:"players"`
	MetaData struct {
		Map              string  `json:"map"`
		GameVersion      string  `json:"game_version"`
		GameStartPatched string  `json:"game_start_patched"`
		Mode             string  `json:"mode"`
		Queue            string  `json:"queue"`
		SeasonID         string  `json:"season_id"`
		Platform         string  `json:"platform"`
		MatchID          string  `json:"match_id"`
		Region           string  `json:"region"`
		Cluster          string  `json:"cluster"`
		GameLength       float64 `json:"game_length"`
		GameStart        float64 `json:"game_start"`
		RoundsPlayed     float64 `json:"rounds_played"`
	}
}

type ValorantMmrHistoryMatch struct {
	Images struct {
		Small        string `json:"small"`
		Large        string `json:"large"`
		TriangleDown string `json:"triangle_down"`
		TriangleUp   string `json:"triangle_up"`
	} `json:"images"`
	Map struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	} `json:"map"`
	CurrenttierPatched  string  `json:"currenttier_patched"`
	MatchId             string  `json:"match_id"`
	SeasonId            string  `json:"season_id"`
	Date                string  `json:"date"`
	Currenttier         float64 `json:"currenttier"`
	RankingInTier       float64 `json:"ranking_in_tier"`
	MmrChangeToLastGame float64 `json:"mmr_change_to_last_game"`
	Elo                 float64 `json:"elo"`
	DateRaw             float64 `json:"date_raw"`
}

type ValorantMmrHistoryResponse struct {
	Data []*ValorantMmrHistoryMatch `json:"data"`
}
