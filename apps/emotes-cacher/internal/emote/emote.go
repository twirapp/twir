package emote

type ID string

type Emote struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`
}
