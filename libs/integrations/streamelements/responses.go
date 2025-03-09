package streamelements

import (
	"time"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type UserProfile struct {
	ID           string `json:"_id"`
	Username     string `json:"username"`
	DisplayName  string `json:"displayName"`
	Avatar       string `json:"avatar"`
	ProfileImage string `json:"profileImage"`
	Provider     string `json:"provider"` // typically "twitch"
}

type CommandCooldown struct {
	User   int `json:"user"`
	Global int `json:"global"`
}

type Command struct {
	ID             string          `json:"_id"`
	Name           string          `json:"command"`
	Enabled        bool            `json:"enabled"`
	Cost           int             `json:"cost"`
	Cooldown       CommandCooldown `json:"cooldown"`
	Aliases        []string        `json:"aliases"`
	Keywords       []string        `json:"keywords"`
	Response       string          `json:"response"`
	AccessLevel    int             `json:"accessLevel"`
	EnabledOnline  bool            `json:"enabledOnline"`
	EnabledOffline bool            `json:"enabledOffline"`
	Hidden         bool            `json:"hidden"`
	Type           string          `json:"type"` // "say", "action", or "custom"
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

type Timer struct {
	Online struct {
		Enabled  bool `json:"enabled"`
		Interval int  `json:"interval"`
	} `json:"online"`
	Offline struct {
		Enabled  bool `json:"enabled"`
		Interval int  `json:"interval"`
	} `json:"offline"`
	Enabled   bool      `json:"enabled"`
	ChatLines int       `json:"chatLines"`
	Id        string    `json:"_id"`
	Channel   string    `json:"channel"`
	Name      string    `json:"name"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
