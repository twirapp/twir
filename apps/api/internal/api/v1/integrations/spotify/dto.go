package spotify

type spotifyDto struct {
	Enabled *bool `validate:"required" json:"enabled"`
}

type tokenDto struct {
	Code string `validate:"required" json:"code"`
}
