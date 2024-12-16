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
