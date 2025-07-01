package ytsr

const SearchSubject = "ytsr.search"

type SearchRequest struct {
	Search    string `json:"search"`
	OnlyLinks bool   `json:"only_links"`
}

type Song struct {
	Title        string     `json:"title"`
	Id           string     `json:"id"`
	Views        uint64     `json:"views"`
	Duration     uint64     `json:"duration"`
	ThumbnailUrl *string    `json:"thumbnail_url"`
	IsLive       bool       `json:"is_live"`
	Author       SongAuthor `json:"author"`
	Link         *string    `json:"link"`
	AuthorName   string     `json:"author_name"`
	AuthorId     string     `json:"author_id"`
	AuthorImage  string     `json:"author_image"`
}

type SongAuthor struct {
	Name      string  `json:"name"`
	ChannelId string  `json:"channel_id"`
	AvatarUrl *string `json:"avatar_url"`
}

type SearchResponse struct {
	Songs []Song `json:"songs"`
}
