package model

type Integration struct {
	ID           string  `db:"id"`
	Service      Service `db:"service"`
	AccessToken  *string `db:"accessToken"`
	RefreshToken *string `db:"refreshToken"`
	ClientID     *string `db:"clientId"`
	ClientSecret *string `db:"clientSecret"`
	APIKey       *string `db:"apiKey"`
	RedirectURL  *string `db:"redirectUrl"`
}

type Service string

const (
	ServiceLastfm         Service = "LASTFM"
	ServiceVK             Service = "VK"
	ServiceFaceit         Service = "FACEIT"
	ServiceSpotify        Service = "SPOTIFY"
	ServiceDonationAlerts Service = "DONATIONALERTS"
	ServiceDiscord        Service = "DISCORD"
	ServiceStreamLabs     Service = "STREAMLABS"
	ServiceDonatePay      Service = "DONATEPAY"
	ServiceDonatello      Service = "DONATELLO"
	ServiceValorant       Service = "VALORANT"
	ServiceDonateStream   Service = "DONATE_STREAM"
	ServiceNightbot       Service = "NIGHTBOT"
)

var Nil = Integration{}
