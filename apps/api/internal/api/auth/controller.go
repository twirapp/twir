package auth

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	cfg "github.com/satont/tsuwari/libs/config"
	"net/http"
	"time"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

var scopes = []string{
	"moderation:read",
	"channel:manage:broadcast",
	"channel:read:redemptions",
	"moderator:read:chatters",
	"moderator:manage:shoutouts",
	"moderator:manage:banned_users",
	"channel:read:vips",
	"channel:manage:vips",
	"channel:manage:moderators",
	"moderator:read:followers",
	"moderator:manage:chat_settings",
}

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("auth")
	middleware.Get("", get())
	middleware.Get("token", getTokens(services))
	middleware.Post("token", refreshToken(services))
	middleware.Post("logout", middlewares.CheckUserAuth(services), logout(services))

	profileCache := cache.New(cache.Config{
		Expiration: 24 * time.Hour,
		Storage:    services.RedisStorage,
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf("fiber:cache:auth:profile:%s", c.Locals("dbUser").(model.Users).ID)
		},
	})
	middleware.Get(
		"profile",
		checkScopes,
		middlewares.CheckUserAuth(services),
		profileCache,
		getProfile(services),
	)

	return middleware
}

func get() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		config := do.MustInvoke[cfg.Config](di.Provider)

		twitchClient, err := helix.NewClient(&helix.Options{
			ClientID:    config.TwitchClientId,
			RedirectURI: config.TwitchCallbackUrl,
		})

		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, "internal error")
		}

		state := c.Query("state")
		if state == "" {
			return c.JSON(fiber.Map{"message": "state is missed"})
		}

		url := twitchClient.GetAuthorizationURL(&helix.AuthorizationURLParams{
			ResponseType: "code",
			Scopes:       scopes,
			State:        state,
			ForceVerify:  false,
		})

		return c.Redirect(url)
	}
}

func getTokens(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		code := c.Query("code")
		state := c.Query("state")

		if code == "" || state == "" {
			return c.Status(401).JSON(fiber.Map{"message": "code or state is missed in request"})
		}

		tokens, err := handleGetToken(code, services)
		if err != nil {
			return err
		}

		services.RedisStorage.DeleteByMethod(
			fmt.Sprintf("fiber:cache:auth:profile:%s", tokens.UserId),
			"GET",
		)

		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    tokens.RefreshToken,
			HTTPOnly: true,
			Expires:  time.Now().UTC().Add(refreshLifeTime),
			SameSite: "lax",
		})
		return c.JSON(fiber.Map{"accessToken": tokens.AccessToken})
	}
}

func getProfile(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Locals("dbUser")
		if user == nil {
			return fiber.NewError(401, "unauthorized")
		}

		profile, err := handleGetProfile(user.(model.Users), services)
		if err != nil {
			return err
		}
		return c.JSON(profile)
	}
}

func logout(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("dbUser").(model.Users).ID

		services.RedisStorage.DeleteByMethod(
			fmt.Sprintf("fiber:cache:auth:profile:%s", userId),
			"GET",
		)

		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    "",
			HTTPOnly: true,
			Expires:  time.Now().UTC(),
			SameSite: "lax",
		})
		return c.SendStatus(200)
	}
}

func refreshToken(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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
}
