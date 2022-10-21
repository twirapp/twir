package commands

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/satont/tsuwari/apps/api-go/internal/middlewares"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("commands")

	middleware.Get(
		"",
		cache.New(cache.Config{
			Expiration: 10 * time.Second,
			Storage:    services.RedisStorage,
			KeyGenerator: func(c *fiber.Ctx) string {
				return fmt.Sprintf("channels:commandsList:%s", c.Params("channelId"))
			},
		}),
		func(c *fiber.Ctx) error {
			c.JSON(HandleGet(c.Params("channelId"), services))

			return nil
		})

	middleware.Post("", func(c *fiber.Ctx) error {
		dto := &commandDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		cmd, err := HandlePost(c.Params("channelId"), services, dto)
		if err == nil {
			c.JSON(cmd)
			return nil
		}

		return err
	})
	middleware.Delete(":commandId", func(c *fiber.Ctx) error {
		err := HandleDelete(c.Params("channelId"), c.Params("commandId"), services)
		if err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusOK)
	})

	return router
}
