package auth_handlers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/tsuwari/apps/api-new/internal/http/helpers"
	"github.com/satont/tsuwari/libs/crypto"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/eventsub"
	"github.com/satont/tsuwari/libs/grpc/generated/scheduler"
	"github.com/satont/tsuwari/libs/twitch"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

type postCodeDto struct {
	Code string `validate:"required" json:"code"`
}

type SessionUser struct {
	helix.User

	ApiKey     string `json:"apiKey"`
	IsBotAdmin bool   `json:"isBotAdmin"`
}

func (c *AuthHandlers) PostCode(ctx *fiber.Ctx) error {
	body := &postCodeDto{}
	if err := c.middlewares.ValidateBody(ctx, body); err != nil {
		return err
	}

	twitchClient, err := twitch.NewAppClient(*c.config, c.grpcClients.Tokens)
	if err != nil {
		c.logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	tokenResponse, err := twitchClient.RequestUserAccessToken(body.Code)
	if err != nil {
		c.logger.Error(err)
		return fiber.NewError(http.StatusUnauthorized, "cannot get user tokens")
	}

	if tokenResponse.ErrorMessage != "" {
		c.logger.Error(tokenResponse.ErrorMessage)
		return fiber.NewError(http.StatusUnauthorized, "wrong code used")
	}

	twitchClient.SetUserAccessToken(tokenResponse.Data.AccessToken)

	users, err := twitchClient.GetUsers(&helix.UsersParams{})
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized, "cannot get user tokens")
	}

	if len(users.Data.Users) == 0 {
		return helpers.CreateBusinessErrorWithMessage(http.StatusInternalServerError, "no user found")
	}

	user := users.Data.Users[0]

	session, err := c.sessionStorage.Get(ctx)
	if err != nil {
		c.logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	dbUser, err := c.getOrCreateUser(user.ID, tokenResponse)
	if err != nil {
		c.logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	sessionUser := SessionUser{
		User:       user,
		ApiKey:     dbUser.ApiKey,
		IsBotAdmin: dbUser.IsBotAdmin,
	}

	session.Set("user", sessionUser)
	if err = session.Save(); err != nil {
		c.logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	c.cacheStorage.DeleteGet("/auth/profile")

	return ctx.SendStatus(http.StatusOK)
}

func (c *AuthHandlers) getOrCreateUser(userId string, token *helix.UserAccessTokenResponse) (*model.Users, error) {
	user := &model.Users{}
	err := c.gorm.
		Where("id = ?", userId).
		Preload("Channel").
		Preload("Token").
		Find(user).Error
	if err != nil {
		return nil, err
	}

	accessToken, err := crypto.Encrypt(token.Data.AccessToken, c.config.TokensCipherKey)
	if err != nil {
		return nil, err
	}

	refreshToken, err := crypto.Encrypt(token.Data.RefreshToken, c.config.TokensCipherKey)
	if err != nil {
		return nil, err
	}

	defaultBot := &model.Bots{}
	err = c.gorm.Where("type = ?", "DEFAULT").Select("id").Find(defaultBot).Error
	if err != nil {
		return nil, err
	}
	if defaultBot.ID == "" {
		c.logger.Error("no default bot found")
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if user.ID == "" {
		user.ID = userId
		user.ApiKey = uuid.NewV4().String()
		user.Token = &model.Tokens{
			ID:                  uuid.NewV4().String(),
			AccessToken:         accessToken,
			RefreshToken:        refreshToken,
			ExpiresIn:           int32(token.Data.ExpiresIn),
			ObtainmentTimestamp: time.Now().UTC(),
			Scopes:              token.Data.Scopes,
		}
		user.Channel = &model.Channels{
			ID:        uuid.NewV4().String(),
			IsEnabled: false,
			BotID:     defaultBot.ID,
		}
		err = c.gorm.Create(user).Error
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	if user.Channel == nil {
		user.Channel = &model.Channels{
			ID:        uuid.NewV4().String(),
			IsEnabled: false,
			BotID:     defaultBot.ID,
		}
		err = c.gorm.Save(user.Channel).Error
		if err != nil {
			return nil, err
		}
	}

	if user.Token == nil {
		user.Token = &model.Tokens{
			ID:                  uuid.NewV4().String(),
			AccessToken:         accessToken,
			RefreshToken:        refreshToken,
			ExpiresIn:           int32(token.Data.ExpiresIn),
			ObtainmentTimestamp: time.Now().UTC(),
			Scopes:              token.Data.Scopes,
		}
		err = c.gorm.Create(user.Token).Error
	} else {
		fmt.Println("updating tokens")
		user.Token.AccessToken = accessToken
		user.Token.RefreshToken = refreshToken
		user.Token.ExpiresIn = int32(token.Data.ExpiresIn)
		user.Token.ObtainmentTimestamp = time.Now().UTC()
		user.Token.Scopes = token.Data.Scopes
		err = c.gorm.Save(user.Token).Error
		if err != nil {
			return nil, err
		}
	}

	_, err = c.grpcClients.EventSub.SubscribeToEvents(
		context.Background(),
		&eventsub.SubscribeToEventsRequest{
			ChannelId: userId,
		},
	)
	if err != nil {
		c.logger.Error(err)
	}

	_, err = c.grpcClients.Scheduler.CreateDefaultCommandsAndRoles(
		context.Background(),
		&scheduler.CreateDefaultCommandsAndRolesRequest{
			UserId: userId,
		},
	)
	if err != nil {
		c.logger.Error(err)
	}

	return user, nil
}
