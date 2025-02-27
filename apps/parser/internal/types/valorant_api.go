package types

type ValorantProfileImages struct {
	Small        string `json:"small"`
	Large        string `json:"large"`
	TriangleDown string `json:"triangle_down"`
	TriangleUp   string `json:"triangle_up"`
}

type ValorantProfile struct {
	Data struct {
		BySeason struct {
			E1A2 struct {
				Error string `json:"error"`
			} `json:"e1a2"`
			E1A3 struct {
				Error string `json:"error"`
			} `json:"e1a3"`
			E2A2 struct {
				Error string `json:"error"`
			} `json:"e2a2"`
			E2A3 struct {
				Error string `json:"error"`
			} `json:"e2a3"`
			E3A1 struct {
				Error string `json:"error"`
			} `json:"e3a1"`
			E3A2 struct {
				Error string `json:"error"`
			} `json:"e3a2"`
			E3A3 struct {
				Error string `json:"error"`
			} `json:"e3a3"`
			E4A1 struct {
				Error string `json:"error"`
			} `json:"e4a1"`
			E4A2 struct {
				Error string `json:"error"`
			} `json:"e4a2"`
			E5A1 struct {
				Error string `json:"error"`
			} `json:"e5a1"`
			E5A2 struct {
				Error string `json:"error"`
			} `json:"e5a2"`
			E5A3 struct {
				Error string `json:"error"`
			} `json:"e5a3"`
			E6A1 struct {
				Error string `json:"error"`
			} `json:"e6a1"`
			E6A2 struct {
				Error string `json:"error"`
			} `json:"e6a2"`
			E6A3 struct {
				Error string `json:"error"`
			} `json:"e6a3"`
			E7A1 struct {
				Error string `json:"error"`
			} `json:"e7a1"`
			E7A2 struct {
				Error string `json:"error"`
			} `json:"e7a2"`
			E7A3 struct {
				Error string `json:"error"`
			} `json:"e7a3"`
			E8A1 struct {
				Error string `json:"error"`
			} `json:"e8a1"`
			E8A2 struct {
				Error string `json:"error"`
			} `json:"e8a2"`
			E8A3 struct {
				Error string `json:"error"`
			} `json:"e8a3"`
			E1A1 struct {
				FinalRankPatched string `json:"final_rank_patched"`
				ActRankWins      []struct {
					PatchedTier string `json:"patched_tier"`
					Tier        int    `json:"tier"`
				} `json:"act_rank_wins"`
				Wins          int  `json:"wins"`
				NumberOfGames int  `json:"number_of_games"`
				FinalRank     int  `json:"final_rank"`
				Old           bool `json:"old"`
			} `json:"e1a1"`
			E2A1 struct {
				FinalRankPatched string `json:"final_rank_patched"`
				ActRankWins      []struct {
					PatchedTier string `json:"patched_tier"`
					Tier        int    `json:"tier"`
				} `json:"act_rank_wins"`
				Wins          int  `json:"wins"`
				NumberOfGames int  `json:"number_of_games"`
				FinalRank     int  `json:"final_rank"`
				Old           bool `json:"old"`
			} `json:"e2a1"`
			E4A3 struct {
				FinalRankPatched string `json:"final_rank_patched"`
				ActRankWins      []struct {
					PatchedTier string `json:"patched_tier"`
					Tier        int    `json:"tier"`
				} `json:"act_rank_wins"`
				Wins          int  `json:"wins"`
				NumberOfGames int  `json:"number_of_games"`
				FinalRank     int  `json:"final_rank"`
				Old           bool `json:"old"`
			} `json:"e4a3"`
		} `json:"by_season"`
		Name        string `json:"name"`
		Tag         string `json:"tag"`
		HighestRank struct {
			PatchedTier string `json:"patched_tier"`
			Season      string `json:"season"`
			Tier        int    `json:"tier"`
			Converted   int    `json:"converted"`
			Old         bool   `json:"old"`
		} `json:"highest_rank"`
		CurrentData struct {
			Images struct {
				Small        string `json:"small"`
				Large        string `json:"large"`
				TriangleDown string `json:"triangle_down"`
				TriangleUp   string `json:"triangle_up"`
			} `json:"images"`
			Currenttierpatched   string `json:"currenttierpatched"`
			Currenttier          int    `json:"currenttier"`
			RankingInTier        int    `json:"ranking_in_tier"`
			MmrChangeToLastGame  int    `json:"mmr_change_to_last_game"`
			Elo                  int    `json:"elo"`
			GamesNeededForRating int    `json:"games_needed_for_rating"`
			Old                  bool   `json:"old"`
		} `json:"current_data"`
	} `json:"data"`
	Status int `json:"status"`
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
