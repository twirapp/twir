package streamlabs

import (
	"fmt"
	"time"
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("streamlabs")
	middleware.Get("auth", getAuth(services))
	middleware.Get("", get(services))

	limit := limiter.New(limiter.Config{
		Max:        1,
		Expiration: 2 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			dbUser := c.Locals("dbUser").(model.Users)
			return fmt.Sprintf("fiber:limiter:integrations:streamlabs:%s", dbUser.ID)
		},
		LimitReached: func(c *fiber.Ctx) error {
			header := c.GetRespHeader("Retry-After", "2")
			return c.Status(429).JSON(fiber.Map{"message": fmt.Sprintf("wait %s seconds", header)})
		},
		Storage: services.RedisStorage,
	})

	middleware.Post("token", post((services)))
	middleware.Patch("", limit, patch((services)))

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
		dto := &streamlabsDto{}
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

type tokenDto struct {
	Code string `validate:"required" json:"code"`
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

		err = handlePost(c.Params("channelId"), dto, services)
		if err != nil {
			return err
		}

		return c.SendStatus(200)
	}
}
