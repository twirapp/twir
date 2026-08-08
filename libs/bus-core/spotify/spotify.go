package spotify

const (
	SearchSubject            = "spotify.search"
	CreateSongRequestSubject = "spotify.songRequest.create"
	CancelSongRequestSubject = "spotify.songRequest.cancel"
)

type TrackData struct {
	ID         string
	URI        string
	Name       string
	ArtistName string
	AlbumName  string
	DurationMs int
	ImageURL   string
}

type SearchRequest struct {
	ChannelID string
	Query     string
	Limit     int
}

type SearchResponse struct {
	Tracks []TrackData
}

type SongRequestData struct {
	ID                   string
	ChannelID            string
	TrackID              string
	TrackURI             string
	Title                string
	Artist               string
	Album                string
	DurationMs           int
	RequesterUserID      string
	RequesterName        string
	RequesterDisplayName string
	Source               string
	QueuePosition        int
	Status               string
	CreatedAt            string
}

type CreateSongRequestRequest struct {
	ChannelID            string
	RequesterUserID      string
	RequesterName        string
	RequesterDisplayName string
	Source               string
	Query                string
}

type CreateSongRequestResponse struct {
	Request SongRequestData
}

type CancelSongRequestRequest struct {
	ChannelID     string
	RequesterName string
}

type CancelSongRequestResponse struct {
	Request SongRequestData
}
