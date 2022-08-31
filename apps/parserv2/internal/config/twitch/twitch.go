package twitch

import (
	"fmt"
	"tsuwari/parser/internal/config/cfg"

	helix "tsuwari/parser/internal/helix"
)

func New(cfg cfg.Config) *helix.Client {
	options := &helix.Options{
		ClientID:     cfg.TwitchClientId,
		ClientSecret: cfg.TwitchClientSecret,
	}
	client, err := helix.NewClient(options)

	if err != nil {
		panic(err)
	}

	token, err := client.RequestAppAccessToken([]string{})

	if err != nil {
		panic(err)
	}

	client.SetAppAccessToken(token.Data.AccessToken, &token.Data.ExpiresIn)

	fmt.Println(client.GetAppAccessToken())

	return client
}
