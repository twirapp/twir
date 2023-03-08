package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/libs/crypto"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/eventsub"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var checkScopes = func(services *types.Services) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		headers := ctx.GetReqHeaders()
		header, ok := headers["Authorization"]
		_, okApiKey := headers["Api-Key"]
		if okApiKey {
			return ctx.Next()
		}

		if !ok {
			return fiber.NewError(http.StatusUnauthorized, "no token provided")
		}
		token, err := middlewares.ExtractTokenFromHeader(services, header)
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
}

// used for register user in system
func checkUser(
	userId string,
	tokens helix.AccessCredentials,
	services *types.Services,
) error {
	defaultBot := model.Bots{}
	err := services.Gorm.Where("type = ?", "DEFAULT").First(&defaultBot).Error
	if err != nil {
		return errors.New("bot not created, cannot create user")
	}

	accessToken, err := crypto.Encrypt(tokens.AccessToken, services.Config.TokensCipherKey)
	if err != nil {
		return err
	}

	refreshToken, err := crypto.Encrypt(tokens.RefreshToken, services.Config.TokensCipherKey)
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
	err = services.Gorm.
		Where(`"users"."id" = ?`, userId).
		Joins("Channel").
		Joins("Token").
		First(&user).Error

	if err != nil && err == gorm.ErrRecordNotFound {
		err = services.Gorm.Transaction(func(tx *gorm.DB) error {
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
			err = createRolesAndCommand(tx, services, userId)
			if err != nil {
				return err
			}

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

			if err = services.Gorm.Create(&channel).Error; err != nil {
				return err
			}
			user.Channel = &channel
		}

		if user.TokenID.Valid {
			tokenData.ID = user.TokenID.String
			if err = services.Gorm.Select("*").Save(&tokenData).Error; err != nil {
				return err
			}
		} else {
			tokenData.ID = uuid.NewV4().String()
			if err = services.Gorm.Save(&tokenData).Error; err != nil {
				return err
			}
			user.TokenID = sql.NullString{String: tokenData.ID, Valid: true}
			if err := services.Gorm.Save(&user).Error; err != nil {
				services.Logger.Error(err)
			}
		}
	}

	services.Grpc.EventSub.SubscribeToEvents(
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
