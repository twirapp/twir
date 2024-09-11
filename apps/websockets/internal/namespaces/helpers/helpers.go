package helpers

import (
	"errors"

	"github.com/olahol/melody"
	model "github.com/satont/twir/libs/gomodels"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("no user found")

func CheckUserByApiKey(db *gorm.DB, session *melody.Session) error {
	apiKey := session.Request.URL.Query().Get("apiKey")
	if apiKey == "" {
		session.Close()
		return errors.New("no api key")
	}

	dbUser := &model.Users{}
	err := db.Where(`"apiKey" = ?`, apiKey).First(dbUser).Error
	if err != nil {
		zap.S().Errorf(apiKey, err)
		session.Close()
		return ErrUserNotFound
	}

	session.Set("userId", dbUser.ID)

	return nil
}
