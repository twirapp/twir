package twitch

import (
	"time"

	helix "github.com/satont/go-helix/v2"
)

type Token struct {
	tokenExpiresIn *int
	tokenCreatedAt *int64
}

type Twitch struct {
	Token
	Client *helix.Client
}

func New(clientId string, clientSecret string) *Twitch {
	options := &helix.Options{
		ClientID:     clientId,
		ClientSecret: clientSecret,
	}
	client, err := helix.NewClient(options)
	if err != nil {
		panic(err)
	}

	twitch := Twitch{
		Client: client,
		Token:  Token{},
	}

	go func() {
		for {
			if !twitch.isTokenValid() {
				exp := twitch.Refresh()
				twitch.setExpiresAndCreated(exp)
			}
			time.Sleep(time.Duration(*twitch.Token.tokenExpiresIn) * time.Second)
		}
	}()

	return &twitch
}

func (c *Twitch) setExpiresAndCreated(expiresIn int) {
	exp := expiresIn
	c.Token.tokenExpiresIn = &exp
	t := time.Now().UnixMilli()
	c.Token.tokenCreatedAt = &t
}

func (c *Twitch) isTokenValid() bool {
	if c.tokenCreatedAt == nil || c.tokenExpiresIn == nil {
		return false
	}

	curr := time.Now().UnixMilli()
	isExpired := curr > (*c.Token.tokenCreatedAt + int64(*c.Token.tokenExpiresIn))

	return isExpired
}

func (c *Twitch) Refresh() int {
	token, err := c.Client.RequestAppAccessToken([]string{})
	if err != nil {
		panic(err)
	}

	c.Client.SetAppAccessToken(token.Data.AccessToken)
	return token.Data.ExpiresIn
}
