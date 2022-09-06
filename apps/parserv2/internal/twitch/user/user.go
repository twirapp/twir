package usersauth

import (
	"fmt"
	model "tsuwari/parser/internal/models"

	helix "github.com/satont/go-helix"
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

func New(opts UsersServiceOpts) *UsersTokensService {
	service := &UsersTokensService{
		db:           opts.Db,
		clientId:     opts.ClientId,
		clientSecret: opts.ClientSecret,
	}

	return service
}

func (c UsersTokensService) Create(userId string) (*helix.Client, error) {
	user := model.UserWitchToken{}

	err := c.db.Where(`"id" = ?`, userId).Find(&user).Error

	if err != nil {
		return nil, err
	}

	refreshFunc := func(tokenData helix.RefreshTokenResponse) {
		err := c.db.Where(`"id" = ?`, user.Token.ID).Updates(model.Tokens{
			AccessToken:  tokenData.Data.AccessToken,
			RefreshToken: tokenData.Data.RefreshToken,
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
	})

	if err != nil {
		return nil, err
	}

	return client, nil
}
