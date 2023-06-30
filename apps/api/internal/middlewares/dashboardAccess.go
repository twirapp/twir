package middlewares

import (
	"net/http"

	model "github.com/satont/twir/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

var CheckHasAccessToDashboard = func(c *fiber.Ctx) error {
	if c.Locals("dbUser") == nil {
		return fiber.NewError(401, "unauthentificated")
	}
	dbUser := c.Locals("dbUser").(model.Users)

	if dbUser.Channel != nil && dbUser.Channel.IsBanned {
		return c.Status(http.StatusForbidden).
			JSON(fiber.Map{"message": "uvy Jokerge"})
	}

	if dbUser.IsBotAdmin {
		return c.Next()
	}

	channelId := c.Params("channelId")
	if channelId == "" {
		return c.Status(http.StatusInternalServerError).
			JSON(fiber.Map{"message": "channelId not passed. This is probably internal error"})
	}
	if dbUser.ID == channelId {
		return c.Next()
	}

	_, ok := lo.Find(
		dbUser.Roles, func(a model.ChannelRoleUser) bool {
			return a.Role.ChannelID == channelId
		},
	)

	if ok {
		return c.Next()
	}

	return c.JSON(fiber.Map{"message": "you have no access to that channel dashboard"})
}
