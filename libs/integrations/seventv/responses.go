package seventv

type ProfileResponse struct {
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
		Emotes []Emote `json:"emotes"`
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

type Emote struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Flags     int    `json:"flags"`
	Timestamp int64  `json:"timestamp"`
	ActorId   string `json:"actor_id"`
	Data      struct {
		Id        string   `json:"id"`
		Name      string   `json:"name"`
		Flags     int      `json:"flags"`
		Tags      []string `json:"tags"`
		Lifecycle int      `json:"lifecycle"`
		State     []string `json:"state"`
		Listed    bool     `json:"listed"`
		Animated  bool     `json:"animated"`
		Owner     struct {
			Id          string `json:"id"`
			Username    string `json:"username"`
			DisplayName string `json:"display_name"`
			AvatarUrl   string `json:"avatar_url"`
			Style       struct {
			} `json:"style"`
			Roles []string `json:"roles"`
		} `json:"owner"`
		Host struct {
			Url   string `json:"url"`
			Files []struct {
				Name       string `json:"name"`
				StaticName string `json:"static_name"`
				Width      int    `json:"width"`
				Height     int    `json:"height"`
				FrameCount int    `json:"frame_count"`
				Size       int    `json:"size"`
				Format     string `json:"format"`
			} `json:"files"`
		} `json:"host"`
	} `json:"data"`
}
