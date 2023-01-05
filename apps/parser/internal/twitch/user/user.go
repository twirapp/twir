package users_twitch_auth

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	cfg "github.com/satont/tsuwari/libs/config"
	"time"

	model "github.com/satont/tsuwari/libs/gomodels"

	helix "github.com/satont/go-helix/v2"
	"gorm.io/gorm"
)

type UsersTokensService struct {
	db           gorm.DB
	clientId     string
	clientSecret string
}

func New() *UsersTokensService {
	config := do.MustInvoke[cfg.Config](di.Provider)
	db := do.MustInvoke[gorm.DB](di.Provider)

	service := &UsersTokensService{
		db:           db,
		clientId:     config.TwitchClientId,
		clientSecret: config.TwitchClientSecret,
	}

	return service
}

func (c UsersTokensService) Create(userId string) (*helix.Client, error) {
	user := model.UserWitchToken{}

	err := c.db.Where(`id = ?`, userId).Preload("Token").Find(&user).Error
	if err != nil {
		return nil, err
	}

	refreshFunc := func(tokenData helix.RefreshTokenResponse) {
		err := c.db.Where(`"id" = ?`, user.Token.ID).Updates(&model.Tokens{
			AccessToken:         tokenData.Data.AccessToken,
			RefreshToken:        tokenData.Data.RefreshToken,
			ExpiresIn:           int32(tokenData.Data.ExpiresIn),
			ObtainmentTimestamp: time.Now(),
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
