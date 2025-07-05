package tokens

const (
	RequestAppTokenSubject  = "tokens.request_app_token"
	RequestUserTokenSubject = "tokens.request_user_token"
	RequestBotTokenSubject  = "tokens.request_bot_token"
)

type GetUserTokenRequest struct {
	UserId string `json:"userId"`
}

type GetBotTokenRequest struct {
	BotId string `json:"botId"`
}

type UpdateTokenRequest struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	ExpiresIn    int64    `json:"expiresIn"`
	Scopes       []string `json:"scopes"`
}

type TokenResponse struct {
	AccessToken string   `json:"access_token"`
	Scopes      []string `json:"scopes"`
	ExpiresIn   int32    `json:"expires_in"`
}
