package middlewares

import (
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

var CheckHasAccessToDashboard = func(c *fiber.Ctx) error {
	if c.Locals("dbUser") == nil {
		return fiber.NewError(401, "unauthentificated")
	}
	dbUser := c.Locals("dbUser").(model.Users)

	if dbUser.IsBotAdmin {
		return c.Next()
	}

	channelId := c.Params("channelId")
	_, ok := lo.Find(dbUser.DashboardAccess, func(a model.ChannelsDashboardAccess) bool {
		return a.ChannelID == channelId
	})

	if ok {
		return c.Next()
	}

	return c.JSON(fiber.Map{"message": "you have no access to that channel dashboard"})
}
