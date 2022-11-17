package bot

import (
	"fmt"
	"time"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("bot")

	isBotModCache := cache.New(cache.Config{
		Expiration: 15 * time.Second,
		Storage:    services.RedisStorage,
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf("channels:isBotMod:%s", c.Params("channelId"))
		},
	})

	middleware.Get("checkmod", isBotModCache, get(services))

	limit := limiter.New(limiter.Config{
		Max:        2,
		Expiration: 1 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			dbUser := c.Locals("dbUser").(model.Users)
			return fmt.Sprintf("fiber:limiter:bot:connection:%s", dbUser.ID)
		},
		LimitReached: func(c *fiber.Ctx) error {
			header := c.GetRespHeader("Retry-After", "2")
			return c.Status(429).JSON(fiber.Map{"message": fmt.Sprintf("wait %s seconds", header)})
		},
		Storage: services.RedisStorage,
	})

	middleware.Patch("connection", limit, patch(services))

	return middleware
}

func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		isBotMod, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}

		if isBotMod == nil {
			return c.JSON(false)
		}

		return c.JSON(*isBotMod)
	}
}

func patch(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &connectionDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		err = handlePatch(c.Params("channelId"), dto, services)

		if err != nil {
			return err
		}

		return c.SendStatus(200)
	}
}
