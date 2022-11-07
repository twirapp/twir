package auth

import (
	"fmt"
	"time"
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("auth")
	middleware.Get("", get(services))
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
		middlewares.CheckUserAuth(services),
		profileCache,
		getProfile(services),
	)

	return middleware
}

var scopes = []string{"moderation:read", "channel:manage:broadcast"}

func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		state := c.Query("state")
		if state == "" {
			return c.JSON(fiber.Map{"message": "state is missed"})
		}

		url := services.Twitch.Client.GetAuthorizationURL(&helix.AuthorizationURLParams{
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
			Expires:  time.Now().Add(refreshLifeTime),
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
			Expires:  time.Now(),
			SameSite: "lax",
		})
		return c.SendStatus(200)
	}
}

type refreshDto struct {
	RefreshToken string `validate:"required" json:"refreshToken"`
}

func refreshToken(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		refreshToken := c.Cookies("refresh_token")
		if refreshToken == "" {
			return fiber.NewError(401, "unauthorized")
		}
		dto := &refreshDto{RefreshToken: refreshToken}

		newAccess, err := handleRefresh(dto, services)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{"accessToken": newAccess})
	}
}
