package model

import (
	"github.com/guregu/null"
)

type User struct {
	ID                string
	TokenID           null.String
	IsTester          bool
	IsBotAdmin        bool
	ApiKey            string
	IsBanned          bool
	HideOnLandingPage bool
}

var Nil = User{}
