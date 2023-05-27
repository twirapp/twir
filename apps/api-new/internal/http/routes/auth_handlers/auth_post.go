package auth_handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/tsuwari/apps/api-new/internal/http/helpers"
	"github.com/satont/tsuwari/libs/twitch"
	"net/http"
)

type postCodeDto struct {
	Code string `validate:"required"`
}

type SessionUser struct {
	helix.User

	ApiKey     string
	IsBotAdmin bool
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

	resp, err := twitchClient.RequestUserAccessToken(body.Code)
	if err != nil {
		c.logger.Error(err)
		return fiber.NewError(http.StatusUnauthorized, "cannot get user tokens")
	}

	twitchClient.SetUserAccessToken(resp.Data.AccessToken)

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

	sessionUser := SessionUser{
		User: user,
	}

	session.Set("user", sessionUser)
	if err = session.Save(); err != nil {
		c.logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	//services.RedisStorage.DeleteByMethod(
	//	fmt.Sprintf("fiber:cache:auth:profile:%s", tokens.UserId),
	//	"GET",
	//)

	return nil
}
