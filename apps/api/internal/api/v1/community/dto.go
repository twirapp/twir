package community

type resetStatsDto struct {
	Field string `validate:"oneof=messages emotes watched" json:"field"`
}
