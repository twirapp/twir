package keywords

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("keywords")

	middleware.Get("", get(services))
	middleware.Post("", post(services))
	middleware.Delete(":keywordId", delete(services))
	middleware.Put(":keywordId", delete(services))

	return middleware
}

func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
	}
}

func post(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
	}
}

func delete(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
	}
}

func put(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
	}
}
