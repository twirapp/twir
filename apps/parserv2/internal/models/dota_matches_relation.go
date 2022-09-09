package model

type DotaMatchWithRelation struct {
	DotaMatches
	GameMode DotaGameModes       `gorm:"foreignKey:GameModeID"`
	Result   *DotaMatchesResults `gorm:"foreignKey:match_id;references:MatchID"`
}
