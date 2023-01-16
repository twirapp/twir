package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("users")

	middleware.Post("ignored", ignoredUsersPost(services))

	return router
}

type ignoredUserPostDto struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type ignoredUsersPostDto struct {
	Users []ignoredUserPostDto `json:"users"`
}

func ignoredUsersPost(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &ignoredUsersPostDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		err = handleIgnoredUsersPost(services, dto)
		if err == nil {
			return c.SendStatus(200)
		}

		return err
	}
}
