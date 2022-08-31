package twitch

import (
	"time"
	"tsuwari/parser/internal/config/cfg"

	helix "github.com/nicklaw5/helix"
)

type Twitch struct {
	tokenExpiresIn *int
	tokenCreatedAt *int64
	Client         *helix.Client
}

func New(cfg cfg.Config) *Twitch {
	options := &helix.Options{
		ClientID:     cfg.TwitchClientId,
		ClientSecret: cfg.TwitchClientSecret,
	}
	client, err := helix.NewClient(options)

	if err != nil {
		panic(err)
	}

	twitch := Twitch{
		Client: client,
	}

	twitch.RefreshIfNeeded()

	return &twitch
}

func (c Twitch) setExpiresAndCreated(expiresIn int) {
	c.tokenExpiresIn = &expiresIn
	t := time.Now().Unix()
	c.tokenCreatedAt = &t
}

func (c Twitch) isTokenValid() bool {
	if c.tokenCreatedAt == nil || c.tokenExpiresIn == nil {
		return false
	}

	curr := time.Now().UnixMilli()
	isExpired := curr > (*c.tokenCreatedAt + int64(*c.tokenExpiresIn))

	return isExpired
}

func (c Twitch) RefreshIfNeeded() {
	if c.isTokenValid() {
		return
	}

	token, err := c.Client.RequestAppAccessToken([]string{})

	if err != nil {
		panic(err)
	}

	c.Client.SetAppAccessToken(token.Data.AccessToken)
	c.setExpiresAndCreated(token.Data.ExpiresIn)
}
