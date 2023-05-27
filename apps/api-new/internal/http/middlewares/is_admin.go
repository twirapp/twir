package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-new/internal/http/helpers"
	model "github.com/satont/tsuwari/libs/gomodels"
	"net/http"
)

func (c *Middlewares) IsAdmin(ctx *fiber.Ctx) error {
	if ctx.Locals("dbUser") == nil {
		return fiber.NewError(401, "unauthentificated")
	}
	dbUser, ok := ctx.Locals("dbUser").(model.Users)

	if !ok {
		c.logger.Error("cannot cast dbUser", dbUser)
		return helpers.CreateBusinessErrorWithMessage(http.StatusInternalServerError, "internal server error")
	}

	if dbUser.IsBotAdmin {
		return ctx.Next()
	} else {
		return ctx.JSON(fiber.Map{"message": "you are not admin"})
	}
}
