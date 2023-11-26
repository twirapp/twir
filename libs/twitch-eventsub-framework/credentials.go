package eventsub_framework

// Credentials represents a method of obtaining Twitch API client credentials
type Credentials interface {
	ClientID() (string, error)
	AppToken() (string, error)
}

type staticCredentials struct {
	id    string
	token string
}

// NewStaticCredentials creates a Credentials instance with a fixed ClientID
// string and AppToken string.
//
// This Credentials implementation should only be used for development as the
// app token will eventually expire and API calls will subsequently fail.
func NewStaticCredentials(clientID string, appToken string) Credentials {
	return &staticCredentials{
		id:    clientID,
		token: appToken,
	}
}

func (s *staticCredentials) ClientID() (string, error) {
	return s.id, nil
}

func (s *staticCredentials) AppToken() (string, error) {
	return s.token, nil
}
