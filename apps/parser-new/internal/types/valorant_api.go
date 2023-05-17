package types

type ValorantProfileImages struct {
	Small        string `json:"small"`
	Large        string `json:"large"`
	TriangleDown string `json:"triangle_down"`
	TriangleUp   string `json:"triangle_up"`
}

type ValorantProfile struct {
	Status int `json:"status"`
	Data   struct {
		CurrentTier        int    `json:"currenttier"`
		CurrentTierpatched string `json:"currenttierpatched"`
		Images             struct {
			Small        string `json:"small"`
			Large        string `json:"large"`
			TriangleDown string `json:"triangle_down"`
			TriangleUp   string `json:"triangle_up"`
		} `json:"images"`
		RankingInTier       int    `json:"ranking_in_tier"`
		MmrChangeToLastGame int    `json:"mmr_change_to_last_game"`
		Elo                 int    `json:"elo"`
		Name                string `json:"name"`
		Tag                 string `json:"tag"`
		Old                 bool   `json:"old"`
	} `json:"data"`
}

type ValorantMatchPlayer struct {
	Name               string `json:"name"`
	Tag                string `json:"tag"`
	Team               string `json:"team"`
	Level              int    `json:"level"`
	Character          string `json:"character"`
	CurrentTier        int    `json:"currenttier"`
	CurrentTierPatched string `json:"currenttier_patched"`
	Behavior           struct {
		AfkRounds    float64 `json:"afk_rounds"`
		FriendlyFire struct {
			Incoming int `json:"incoming"`
			Outgoing int `json:"outgoing"`
		} `json:"friendly_fire"`
		RoundsInSpawn int `json:"rounds_in_spawn"`
	} `json:"behavior"`
	Stats struct {
		Score     int `json:"score"`
		Kills     int `json:"kills"`
		Deaths    int `json:"deaths"`
		Assists   int `json:"assists"`
		Bodyshots int `json:"bodyshots"`
		Headshots int `json:"headshots"`
		Legshots  int `json:"legshots"`
	} `json:"stats"`
	Economy struct {
		Spent struct {
			Overall int `json:"overall"`
			Average int `json:"average"`
		} `json:"spent"`
		LoadoutValue struct {
			Overall int `json:"overall"`
			Average int `json:"average"`
		} `json:"loadout_value"`
	} `json:"economy"`
	DamageMade     int `json:"damage_made"`
	DamageReceived int `json:"damage_received"`
}

type ValorantMatchPlayers struct {
	AllPlayers []ValorantMatchPlayer `json:"all_players"`
}

type ValorantMatchesResponse struct {
	Data []ValorantMatch `json:"data"`
}

type ValorantMatch struct {
	MetaData struct {
		Map              string `json:"map"`
		GameVersion      string `json:"game_version"`
		GameLength       int    `json:"game_length"`
		GameStart        int    `json:"game_start"`
		GameStartPatched string `json:"game_start_patched"`
		RoundsPlayed     int    `json:"rounds_played"`
		Mode             string `json:"mode"`
		Queue            string `json:"queue"`
		SeasonID         string `json:"season_id"`
		Platform         string `json:"platform"`
		MatchID          string `json:"match_id"`
		Region           string `json:"region"`
		Cluster          string `json:"cluster"`
	}
	Players ValorantMatchPlayers `json:"players"`
	Teams   map[string]struct {
		HasWon     bool `json:"has_won"`
		RoundsWon  int  `json:"rounds_won"`
		RoundsLost int  `json:"rounds_lost"`
	}
}
