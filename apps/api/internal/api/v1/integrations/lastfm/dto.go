package lastfm

type lastfmDto struct {
	Code string `validate:"required" json:"code"`
}
