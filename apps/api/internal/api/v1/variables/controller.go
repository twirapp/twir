package variables

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("variables")
	middleware.Get("", get(services))
	middleware.Post("", post(services))
	middleware.Delete(":variableId", delete(services))
	middleware.Put(":variableId", put(services))
	middleware.Get("builtin", getBuiltIn(services))

	return middleware
}

func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		variables, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}
		return c.JSON(variables)
	}
}

func getBuiltIn(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		variables, err := handleGetBuiltIn(services)
		if err != nil {
			return err
		}

		return c.JSON(variables)
	}
}

func post(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &variableDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		variable, err := handlePost(c.Params("channelId"), dto, services)
		if err == nil {
			return c.JSON(variable)
		}

		return err
	}
}

func delete(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := handleDelete(c.Params("channelId"), c.Params("variableId"), services)
		if err != nil {
			return err
		}
		return c.SendStatus(200)
	}
}

func put(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &variableDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		variable, err := handleUpdate(c.Params("channelId"), c.Params("variableId"), dto, services)
		if err == nil {
			return c.JSON(variable)
		}

		return err
	}
}
