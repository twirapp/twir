package commands

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("commands")
	middleware.Get(
		":channelId",
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
	middleware.Post(":channelId", func(c *fiber.Ctx) error {
		dto := commandDto{}
		if err := c.BodyParser(&dto); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "wrong body",
			})
		}

		if err := services.Validator.Struct(dto); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			errors := []string{}
			for _, e := range validationErrors {
				fmt.Println(e.StructField(), e.Translate(services.ValidatorTranslator))
				errors = append(
					errors,
					fmt.Sprintf(
						"%s %s",
						e.StructField(),
						e.Translate(services.ValidatorTranslator),
					),
				)
			}
			c.Status(fiber.StatusBadRequest).JSON(errors)
			return nil
		}

		if err := c.JSON(HandlePost(c.Params("channelId"), services)); err != nil {
			return c.SendString(err.Error())
		}

		return nil
	})

	return router
}
