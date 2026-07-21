package api

const (
	SongRequestAddToQueueSubject      = "api.songRequests.addToQueue"
	SongRequestRemoveFromQueueSubject = "api.songRequests.removeFromQueue"
	SongRequestPlaybackStateSubject   = "api.songRequests.playbackState"
)

type SongRequestAddToQueue struct {
	ChannelID   string
	SongRequest SongRequestData
}

type SongRequestRemoveFromQueue struct {
	ChannelID string
	VideoID   string
}

type SongRequestPlaybackState struct {
	ChannelID string
	VideoID   string
	Title     string
	Position  float64
	IsPlaying bool
	Volume    int
	UpdatedAt int64
}

type SongRequestData struct {
	ID                   string
	Title                string
	VideoID              string
	SongLink             string
	DurationSeconds      int
	OrderedByName        string
	OrderedByDisplayName string
	QueuePosition        int
	CreatedAt            string
}
