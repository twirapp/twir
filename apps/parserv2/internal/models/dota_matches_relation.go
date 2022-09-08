package model

type DotaMatchWithRelation struct {
	DotaMatches
	GameMode DotaGameModes `gorm:"foreignKey:GameModeID"`
}
