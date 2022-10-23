package feedback

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("feedback")

	limit := limiter.New(limiter.Config{
		Max:        2,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			userId := "123123"
			return fmt.Sprintf("fiber:limiter:feedback:%s", userId)
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{"message": "you are sending feedback to fast"})
		},
		Storage: services.RedisStorage,
	})

	middleware.Post("", limit, func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			return fiber.NewError(400, "wrong form data")
		}

		if form.Value["text"] == nil {
			return fiber.NewError(400, "text cannot be empty")
		}

		// TODO: add user id to the params
		err = handlePost("", form.Value["text"][0], form.File["files"], services)
		if err != nil {
			return err
		}

		return c.SendStatus(200)
	})

	return middleware
}
