package messages

type Dispatch struct {
	Type string `json:"type"`
	Body struct {
		Id    string `json:"id"`
		Kind  int    `json:"kind"`
		Actor *struct {
			Id          string `json:"id"`
			Username    string `json:"username"`
			DisplayName string `json:"display_name"`
			AvatarUrl   string `json:"avatar_url"`
			Style       struct {
			} `json:"style"`
			RoleIds     []string `json:"role_ids"`
			Connections []struct {
				Id            string `json:"id"`
				Platform      string `json:"platform"`
				Username      string `json:"username"`
				DisplayName   string `json:"display_name"`
				LinkedAt      int64  `json:"linked_at"`
				EmoteCapacity int    `json:"emote_capacity"`
				EmoteSetId    string `json:"emote_set_id"`
			} `json:"connections"`
		} `json:"actor"`
		Pulled  []DispatchChange `json:"pulled"`  // deletes
		Pushed  []DispatchChange `json:"pushed"`  // creates
		Updated []DispatchChange `json:"updated"` // updates
	} `json:"body"`
}

type DispatchChange struct {
	Key      string              `json:"key"`
	Index    int                 `json:"index"`
	Type     string              `json:"type"`
	OldValue *DispatchEmoteValue `json:"old_value"`
	Value    *DispatchEmoteValue `json:"value"`
}

type DispatchEmoteValue struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Flags     int    `json:"flags"`
	Timestamp int64  `json:"timestamp"`
	ActorId   string `json:"actor_id"`
	Data      *struct {
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
			RoleIds     []string `json:"role_ids"`
			Connections []struct {
				Id            string `json:"id"`
				Platform      string `json:"platform"`
				Username      string `json:"username"`
				DisplayName   string `json:"display_name"`
				LinkedAt      int64  `json:"linked_at"`
				EmoteCapacity int    `json:"emote_capacity"`
				EmoteSetId    string `json:"emote_set_id"`
			} `json:"connections"`
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
	OriginId interface{} `json:"origin_id"`
}
