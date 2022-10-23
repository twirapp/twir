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

	commandsCache := cache.New(cache.Config{
		Expiration: 10 * time.Second,
		Storage:    services.RedisStorage,
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf("channels:commandsList:%s", c.Params("channelId"))
		},
	})

	middleware.Get("", commandsCache, get(services))
	middleware.Post("", post(services))
	middleware.Delete(":commandId", delete(services))
	middleware.Put(":commandId", update(services))

	return router
}

func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.JSON(handleGet(c.Params("channelId"), services))

		return nil
	}
}

func post(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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

		cmd, err := handlePost(c.Params("channelId"), services, dto)
		if err == nil {
			return c.JSON(cmd)
		}

		return err
	}
}

func delete(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := handleDelete(c.Params("channelId"), c.Params("commandId"), services)
		if err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusOK)
	}
}

func update(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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

		cmd, err := handleUpdate(c.Params("channelId"), c.Params("commandId"), dto, services)
		if err == nil && cmd != nil {
			return c.JSON(cmd)
		}

		return err
	}
}
