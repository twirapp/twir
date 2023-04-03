package auth

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/crypto"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/eventsub"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var checkScopes = func(ctx *fiber.Ctx) error {
	headers := ctx.GetReqHeaders()
	header, ok := headers["Authorization"]
	_, okApiKey := headers["Api-Key"]
	if okApiKey {
		return ctx.Next()
	}

	if !ok {
		return fiber.NewError(http.StatusUnauthorized, "no token provided")
	}
	token, err := middlewares.ExtractTokenFromHeader(header)
	if err != nil {
		return ctx.Next()
	}
	claims := token.Claims.(jwt.MapClaims)
	reqScopes, ok := claims["scopes"]
	if !ok {
		return ctx.SendStatus(http.StatusForbidden)
	}

	parsedScopes := lo.Map(reqScopes.([]any), func(item any, _ int) string {
		scope, ok := item.(string)
		if !ok {
			return ""
		}
		return scope
	})

	for _, scope := range scopes {
		_, ok := lo.Find(parsedScopes, func(s string) bool {
			return s == scope
		})

		if !ok {
			return ctx.Status(http.StatusForbidden).SendString("not enough scopes")
		}
	}
	return ctx.Next()
}

// used for register user in system
func checkUser(
	userId string,
	tokens helix.AccessCredentials,
	services types.Services,
) error {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	eventSubGrpc := do.MustInvoke[eventsub.EventSubClient](di.Provider)
	//schedulerGrpc := do.MustInvoke[scheduler.SchedulerClient](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)

	defaultBot := model.Bots{}
	err := services.DB.Where("type = ?", "DEFAULT").First(&defaultBot).Error
	if err != nil {
		return errors.New("bot not created, cannot create user")
	}

	accessToken, err := crypto.Encrypt(tokens.AccessToken, config.TokensCipherKey)
	if err != nil {
		return err
	}

	refreshToken, err := crypto.Encrypt(tokens.RefreshToken, config.TokensCipherKey)
	if err != nil {
		return err
	}

	tokenData := model.Tokens{
		AccessToken:         accessToken,
		RefreshToken:        refreshToken,
		ExpiresIn:           int32(tokens.ExpiresIn),
		ObtainmentTimestamp: time.Now().UTC(),
		Scopes:              tokens.Scopes,
	}

	user := model.Users{}
	err = services.DB.
		Where(`"users"."id" = ?`, userId).
		Joins("Channel").
		Joins("Token").
		First(&user).Error

	if err != nil && err == gorm.ErrRecordNotFound {
		err = services.DB.Transaction(func(tx *gorm.DB) error {
			newToken := tokenData
			newToken.ID = uuid.NewV4().String()

			if err = tx.Save(&newToken).Error; err != nil {
				return err
			}

			user.ID = userId
			user.TokenID = sql.NullString{String: newToken.ID, Valid: true}
			user.ApiKey = uuid.NewV4().String()

			if err = tx.Save(&user).Error; err != nil {
				return err
			}

			channel := createChannelModel(user.ID, defaultBot.ID)

			if err = tx.Create(&channel).Error; err != nil {
				return err
			}
			user.Channel = &channel

			return nil
		})

		if err != nil {
			return err
		}
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	} else {
		if user.Channel == nil {
			channel := createChannelModel(user.ID, defaultBot.ID)

			if err = services.DB.Create(&channel).Error; err != nil {
				return err
			}
			user.Channel = &channel
		}

		if user.TokenID.Valid {
			tokenData.ID = user.TokenID.String
			if err = services.DB.Select("*").Save(&tokenData).Error; err != nil {
				return err
			}
		} else {
			tokenData.ID = uuid.NewV4().String()
			if err = services.DB.Save(&tokenData).Error; err != nil {
				return err
			}
			user.TokenID = sql.NullString{String: tokenData.ID, Valid: true}
			if err := services.DB.Save(&user).Error; err != nil {
				logger.Error(err)
			}
		}
	}

	err = middlewares.CreateRolesAndCommand(services.DB, userId)
	if err != nil {
		return err
	}

	eventSubGrpc.SubscribeToEvents(
		context.Background(),
		&eventsub.SubscribeToEventsRequest{
			ChannelId: userId,
		},
	)

	return nil
}

func createChannelModel(userId, botId string) model.Channels {
	return model.Channels{
		ID:             userId,
		IsEnabled:      false,
		IsTwitchBanned: false,
		IsBanned:       false,
		BotID:          botId,
	}
}
