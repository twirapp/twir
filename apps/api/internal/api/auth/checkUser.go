package auth

import (
	"context"
	"database/sql"
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

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/eventsub"
	"github.com/satont/tsuwari/libs/grpc/generated/scheduler"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
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
	schedulerGrpc := do.MustInvoke[scheduler.SchedulerClient](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)

	defaultBot := model.Bots{}
	err := services.DB.Where("type = ?", "DEFAULT").First(&defaultBot).Error
	if err != nil {
		return fiber.NewError(500, "bot not created, cannot create user")
	}

	accessToken, err := crypto.Encrypt(tokens.AccessToken, config.TokensCipherKey)
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	refreshToken, err := crypto.Encrypt(tokens.RefreshToken, config.TokensCipherKey)
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
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

			if err = services.DB.Save(&newToken).Error; err != nil {
				return err
			}

			user.ID = userId
			user.TokenID = sql.NullString{String: newToken.ID, Valid: true}
			user.ApiKey = uuid.NewV4().String()

			if err = services.DB.Save(&user).Error; err != nil {
				return err
			}

			channel := createChannelModel(user.ID, defaultBot.ID)

			if err = services.DB.Create(&channel).Error; err != nil {
				return err
			}
			user.Channel = &channel
			err = createRolesAndCommand(services, userId)
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			logger.Error(err)
			return fiber.NewError(http.StatusInternalServerError, "internal error")
		}
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(err)
		return fiber.NewError(500, "internal error")
	} else {
		if user.Channel == nil {
			channel := createChannelModel(user.ID, defaultBot.ID)

			if err = services.DB.Create(&channel).Error; err != nil {
				logger.Error(err)
				return err
			}
			user.Channel = &channel
		}

		if user.TokenID.Valid {
			tokenData.ID = user.TokenID.String
			if err = services.DB.Select("*").Save(&tokenData).Error; err != nil {
				logger.Error(err)
				return err
			}
		} else {
			tokenData.ID = uuid.NewV4().String()
			if err = services.DB.Save(&tokenData).Error; err != nil {
				logger.Error(err)
				return err
			}
			user.TokenID = sql.NullString{String: tokenData.ID, Valid: true}
			if err := services.DB.Save(&user).Error; err != nil {
				logger.Error(err)
			}
		}
	}

	schedulerGrpc.CreateDefaultCommands(
		context.Background(),
		&scheduler.CreateDefaultCommandsRequest{
			UserId: userId,
		},
	)
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
