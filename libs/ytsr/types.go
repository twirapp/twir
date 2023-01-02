package youtubego

type Thumbnail struct {
	Id, Url       string
	Width, Height float64
}

type Video struct {
	Thumbnail
	Id, Title, Url string
}

type VideoParser struct {
	Video
	IsSuccess bool
}

type Channel struct {
	Icon                       Thumbnail
	Id, Url, Name, Subscribers string
	Verified                   bool
}

type ChannelParser struct {
	Channel
	IsSuccess bool
}

type Playlist struct {
	Thumbnail
	Channel
	Id, title string
	Videos    int
}

type PlaylistParser struct {
	Playlist
	IsSuccess bool
}

type SearchResult struct {
	Video
	Playlist
	Channel
}

type SearchOptions struct {
	Limit int
	Type  string
}
