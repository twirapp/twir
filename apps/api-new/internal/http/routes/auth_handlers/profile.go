package auth_handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-new/internal/http/helpers"
)

func (c *AuthHandlers) GetProfile(ctx *fiber.Ctx) error {
	session, err := c.sessionStorage.Get(ctx)
	if err != nil {
		c.logger.Error(err)
		return helpers.ErrInternalError
	}

	user, castOk := session.Get("user").(SessionUser)
	if !castOk {
		c.logger.Error("Cannot cast user from session")
		return helpers.ErrInternalError
	}

	return ctx.JSON(user)
}
