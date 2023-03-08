package lastfm

type LastfmProfile struct {
	Name       string `json:"name"`
	Image      string `json:"image"`
	PlayCount  string `json:"playCount"`
	TrackCount string `json:"trackCount"`
	AlbumCount string `json:"albumCount"`
}
