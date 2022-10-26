package auth

import (
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/middlewares"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("auth")
	middleware.Get("", get(services))
	middleware.Get("token", getToken(services))
	middleware.Get("profile", middlewares.CheckUserAuth(services), getProfile(services))

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

func getToken(services types.Services) func(c *fiber.Ctx) error {
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

		return c.JSON(tokens)
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
