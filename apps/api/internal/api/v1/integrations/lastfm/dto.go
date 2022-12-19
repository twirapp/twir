package lastfm

type lastfmDto struct {
	Token string `validate:"required" json:"token"`
}
