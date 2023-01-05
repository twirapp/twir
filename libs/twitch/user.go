package twitch

import (
	"errors"
	"fmt"
	"time"

	model "github.com/satont/tsuwari/libs/gomodels"

	helix "github.com/satont/go-helix/v2"
	"gorm.io/gorm"
)

type Opts struct {
	ClientId     string
	ClientSecret string
}

type UsersServiceOpts struct {
	Db           *gorm.DB
	ClientId     string
	ClientSecret string
}

type UsersTokensService struct {
	db           *gorm.DB
	clientId     string
	clientSecret string
}

func NewUserClient(opts UsersServiceOpts) *UsersTokensService {
	service := &UsersTokensService{
		db:           opts.Db,
		clientId:     opts.ClientId,
		clientSecret: opts.ClientSecret,
	}

	return service
}

func (c UsersTokensService) Create(userId string) (*helix.Client, error) {
	user := model.Users{}

	err := c.db.Where(`id = ?`, userId).Preload("Token").First(&user).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, err
	}

	if user.Token == nil {
		return nil, errors.New("cannot find token of user")
	}

	refreshFunc := func(tokenData helix.RefreshTokenResponse) {
		err := c.db.Where(`"id" = ?`, user.Token.ID).Select("*").Updates(&model.Tokens{
			ID:                  user.Token.ID,
			AccessToken:         tokenData.Data.AccessToken,
			RefreshToken:        tokenData.Data.RefreshToken,
			ExpiresIn:           int32(tokenData.Data.ExpiresIn),
			ObtainmentTimestamp: time.Now().UTC(),
		}).Error
		if err != nil {
			fmt.Println(err)
		}
	}

	client, err := helix.NewClient(&helix.Options{
		ClientID:         c.clientId,
		ClientSecret:     c.clientSecret,
		UserAccessToken:  user.Token.AccessToken,
		UserRefreshToken: user.Token.RefreshToken,
		OnRefresh:        &refreshFunc,
		RateLimitFunc:    rateLimitCallback,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c UsersTokensService) CreateBot(botId string) (*helix.Client, error) {
	bot := model.Bots{}

	err := c.db.Where(`id = ?`, botId).Preload("Token").First(&bot).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("bot not found")
	}

	if err != nil {
		return nil, err
	}

	if bot.Token == nil {
		return nil, errors.New("cannot find token of bot")
	}

	refreshFunc := func(tokenData helix.RefreshTokenResponse) {
		err := c.db.Where(`"id" = ?`, bot.Token.ID).Select("*").Updates(&model.Tokens{
			ID:                  bot.Token.ID,
			AccessToken:         tokenData.Data.AccessToken,
			RefreshToken:        tokenData.Data.RefreshToken,
			ExpiresIn:           int32(tokenData.Data.ExpiresIn),
			ObtainmentTimestamp: time.Now().UTC(),
		}).Error
		if err != nil {
			fmt.Println(err)
		}
	}

	client, err := helix.NewClient(&helix.Options{
		ClientID:         c.clientId,
		ClientSecret:     c.clientSecret,
		UserAccessToken:  bot.Token.AccessToken,
		UserRefreshToken: bot.Token.RefreshToken,
		OnRefresh:        &refreshFunc,
		RateLimitFunc:    rateLimitCallback,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}
