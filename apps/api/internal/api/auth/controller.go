package auth

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/services/redis_storage"
	"github.com/satont/tsuwari/libs/twitch"
	"time"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
)

func Setup(router fiber.Router) fiber.Router {
	middleware := router.Group("auth")
	middleware.Get("", get)
	middleware.Get("token", getTokens)
	middleware.Post("token", refreshToken)
	middleware.Post("logout", middlewares.CheckUserAuth, logout)

	profileCache := cache.New(cache.Config{
		Expiration: 24 * time.Hour,
		Storage:    do.MustInvoke[*redis_storage.RedisStorage](di.Injector),
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf("fiber:cache:auth:profile:%s", c.Locals("dbUser").(model.Users).ID)
		},
	})
	middleware.Get(
		"profile",
		middlewares.CheckUserAuth,
		profileCache,
		getProfile,
	)

	return middleware
}

var scopes = []string{"moderation:read", "channel:manage:broadcast"}

var get = func(c *fiber.Ctx) error {
	state := c.Query("state")
	if state == "" {
		return c.JSON(fiber.Map{"message": "state is missed"})
	}

	twitch := do.MustInvoke[*twitch.Twitch](di.Injector)
	url := twitch.Client.GetAuthorizationURL(&helix.AuthorizationURLParams{
		ResponseType: "code",
		Scopes:       scopes,
		State:        state,
		ForceVerify:  false,
	})

	return c.Redirect(url)
}

var getTokens = func(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		return c.Status(401).JSON(fiber.Map{"message": "code or state is missed in request"})
	}

	tokens, err := handleGetToken(code)
	if err != nil {
		return err
	}

	storage := do.MustInvoke[*redis_storage.RedisStorage](di.Injector)
	storage.DeleteByMethod(
		fmt.Sprintf("fiber:cache:auth:profile:%s", tokens.UserId),
		"GET",
	)

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HTTPOnly: true,
		Expires:  time.Now().Add(refreshLifeTime),
		SameSite: "lax",
	})
	return c.JSON(fiber.Map{"accessToken": tokens.AccessToken})
}

var getProfile = func(c *fiber.Ctx) error {
	user := c.Locals("dbUser")
	if user == nil {
		return fiber.NewError(401, "unauthorized")
	}

	profile, err := handleGetProfile(user.(model.Users))
	if err != nil {
		return err
	}
	return c.JSON(profile)
}

var logout = func(c *fiber.Ctx) error {
	userId := c.Locals("dbUser").(model.Users).ID

	storage := do.MustInvoke[*redis_storage.RedisStorage](di.Injector)
	storage.DeleteByMethod(
		fmt.Sprintf("fiber:cache:auth:profile:%s", userId),
		"GET",
	)

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Expires:  time.Now(),
		SameSite: "lax",
	})
	return c.SendStatus(200)
}

var refreshToken = func(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return fiber.NewError(401, "unauthorized")
	}

	newAccess, err := handleRefresh(refreshToken)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"accessToken": newAccess})
}
