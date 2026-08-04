package streamlabs

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type UserProfile struct {
	StreamLabs struct {
		DisplayName string `json:"display_name"`
		ThumbNail   string `json:"thumbnail"`
		ID          int    `json:"id"`
	} `json:"streamlabs"`
}

type SocketTokenResponse struct {
	SocketToken string `json:"socket_token"`
}
