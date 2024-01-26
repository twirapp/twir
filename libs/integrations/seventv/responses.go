package seventv

type SevenTvProfileResponse struct {
	Id            string      `json:"id"`
	Platform      string      `json:"platform"`
	Username      string      `json:"username"`
	DisplayName   string      `json:"display_name"`
	LinkedAt      int64       `json:"linked_at"`
	EmoteCapacity int         `json:"emote_capacity"`
	EmoteSetId    interface{} `json:"emote_set_id"`
	EmoteSet      *struct {
		Id         string        `json:"id"`
		Name       string        `json:"name"`
		Flags      int           `json:"flags"`
		Tags       []interface{} `json:"tags"`
		Immutable  bool          `json:"immutable"`
		Privileged bool          `json:"privileged"`
		Capacity   int           `json:"capacity"`
		Owner      struct {
			Id          string `json:"id"`
			Username    string `json:"username"`
			DisplayName string `json:"display_name"`
			AvatarUrl   string `json:"avatar_url"`
			Style       struct {
			} `json:"style"`
			Roles []string `json:"roles"`
		} `json:"owner"`
	} `json:"emote_set"`
	User struct {
		Id          string `json:"id"`
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		CreatedAt   int64  `json:"created_at"`
		AvatarUrl   string `json:"avatar_url"`
		Style       struct {
		} `json:"style"`
		Editors []struct {
			Id          string `json:"id"`
			Permissions int    `json:"permissions"`
			Visible     bool   `json:"visible"`
			AddedAt     int64  `json:"added_at"`
		} `json:"editors"`
		Roles       []string `json:"roles"`
		Connections []struct {
			Id            string      `json:"id"`
			Platform      string      `json:"platform"`
			Username      string      `json:"username"`
			DisplayName   string      `json:"display_name"`
			LinkedAt      int64       `json:"linked_at"`
			EmoteCapacity int         `json:"emote_capacity"`
			EmoteSetId    interface{} `json:"emote_set_id"`
			EmoteSet      struct {
				Id         string        `json:"id"`
				Name       string        `json:"name"`
				Flags      int           `json:"flags"`
				Tags       []interface{} `json:"tags"`
				Immutable  bool          `json:"immutable"`
				Privileged bool          `json:"privileged"`
				Capacity   int           `json:"capacity"`
				Owner      interface{}   `json:"owner"`
			} `json:"emote_set"`
		} `json:"connections"`
	} `json:"user"`
}
