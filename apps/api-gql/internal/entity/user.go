package entity

type User struct {
	ID                string
	TokenID           *string
	IsTester          bool
	IsBotAdmin        bool
	ApiKey            string
	IsBanned          bool
	HideOnLandingPage bool
}

var UserNil = User{}

type TwitchUser struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"displayName"`
	ProfileImageURL string `json:"profileImageUrl"`
	Description     string `json:"description"`
	NotFound        bool   `json:"notFound"`
}
