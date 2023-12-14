package model

type ChannelModulesSettingsDuel struct {
	Enabled         bool   `json:"enabled"`
	UserCooldown    int32  `json:"user_cooldown"`
	GlobalCooldown  int32  `json:"global_cooldown"`
	TimeoutSeconds  int32  `json:"timeout_seconds"`
	StartMessage    string `json:"start_message"`
	ResultMessage   string `json:"result_message"`
	SecondsToAccept int32  `json:"seconds_to_accept"`
	PointsPerWin    int32  `json:"points_per_win"`
	PointsPerLose   int32  `json:"points_per_lose"`
	BothDiePercent  int32  `json:"both_die_percent"`
	BothDieMessage  string `json:"both_die_message"`
}
