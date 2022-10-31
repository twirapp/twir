package faceit

type faceitData struct {
	UserName string  `validate:"required" json:"username"`
	Game     *string `                    json:"game"`
}

type faceitUpdateDto struct {
	Data    faceitData `validae:"required" json:"data"`
	Enabled *bool      `                   json:"enabled" validate:"required"`
}
