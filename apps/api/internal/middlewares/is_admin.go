package middlewares

import (
	"github.com/gofiber/fiber/v2"
	model "github.com/satont/twir/libs/gomodels"
)

var IsAdmin = func(c *fiber.Ctx) error {
	if c.Locals("dbUser") == nil {
		return fiber.NewError(401, "unauthentificated")
	}
	dbUser := c.Locals("dbUser").(model.Users)

	if dbUser.IsBotAdmin {
		return c.Next()
	} else {
		return c.JSON(fiber.Map{"message": "you are not admin"})
	}
}
