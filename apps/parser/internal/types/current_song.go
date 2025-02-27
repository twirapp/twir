package types

type CurrentSong struct {
	Playlist *CurrentSongPlayList `json:"playlist"`
	Name     string               `json:"song"`
	Image    string               `json:"coverUrl"`
}

type CurrentSongPlayList struct {
	Name      *string `json:"name"`
	Image     *string `json:"coverUrl"`
	Followers *int    `json:"followers"`
	Href      string  `json:"href"`
}
