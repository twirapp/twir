package model

type ChannelGamesDuel struct {
	ID              string `gorm:"column:id;type:uuid"`
	ChannelID       string `gorm:"column:channel_id;type:text"`
	Enabled         bool   `gorm:"column:enabled;type:bool"`
	UserCooldown    int32  `gorm:"column:user_cooldown;type:int"`
	GlobalCooldown  int32  `gorm:"column:global_cooldown;type:int"`
	TimeoutSeconds  int32  `gorm:"column:timeout_seconds;type:int"`
	StartMessage    string `gorm:"column:start_message;type:text"`
	ResultMessage   string `gorm:"column:result_message;type:text"`
	SecondsToAccept int32  `gorm:"column:seconds_to_accept;type:int"`
	PointsPerWin    int32  `gorm:"column:points_per_win;type:int"`
	PointsPerLose   int32  `gorm:"column:points_per_lose;type:int"`
	BothDiePercent  int32  `gorm:"column:both_die_percent;type:int"`
	BothDieMessage  string `gorm:"column:both_die_message;type:text"`
}

func (c ChannelGamesDuel) TableName() string {
	return "channels_games_duel"
}
