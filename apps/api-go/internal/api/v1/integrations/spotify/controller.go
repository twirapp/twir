package spotify

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/satont/tsuwari/apps/api-go/internal/types"

	"github.com/satont/tsuwari/apps/api-go/internal/middlewares"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("spotify")
	middleware.Get("auth", getAuth(services))
	middleware.Get("", get(services))
	middleware.Post("token", post((services)))
	middleware.Patch("", patch((services)))

	profileCache := cache.New(cache.Config{
		Expiration: 31 * 24 * time.Hour,
		Storage:    services.RedisStorage,
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf("fiber:cache:integrations:spotify:profile:%s", c.Params("channelId"))
		},
	})

	middleware.Get("profile", profileCache, getProfile((services)))

	return middleware
}

func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		integration, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}
		return c.JSON(integration)
	}
}

func getProfile(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		profile, err := handleGetProfile(c.Params("channelId"), services)
		if err != nil {
			return err
		}
		return c.JSON(profile)
	}
}

func getAuth(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authLink, err := handleGetAuth(services)
		if err != nil {
			return err
		}

		return c.SendString(*authLink)
	}
}

func patch(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &spotifyDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		integration, err := handlePatch(c.Params("channelId"), dto, services)
		if err != nil {
			return err
		}

		return c.JSON(integration)
	}
}

func post(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &tokenDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		channelId := c.Params("channelId")
		err = handlePost(channelId, dto, services)
		if err != nil {
			return err
		}

		services.RedisStorage.Delete(
			fmt.Sprintf("fiber:cache:integrations:spotify:profile:%s", channelId),
		)

		return c.SendStatus(200)
	}
}
