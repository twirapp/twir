package bot

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
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

	middleware.Get("/checkmod", isBotModCache, get(services))

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
