package lastfm

type lastfmDataDto struct {
	UserName string `validate:"required" json:"username"`
}

type lastfmDto struct {
	Enabled *bool         `validate:"required" json:"enabled"`
	Data    lastfmDataDto `validate:"required" json:"data"`
}
