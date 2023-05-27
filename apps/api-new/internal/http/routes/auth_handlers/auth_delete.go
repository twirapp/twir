package auth_handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-new/internal/http/helpers"
	"net/http"
)

func (c *AuthHandlers) Logout(ctx *fiber.Ctx) error {
	session, err := c.sessionStorage.Get(ctx)
	if err != nil {
		c.logger.Error(err)
		return helpers.ErrInternalError
	}

	if err = session.Destroy(); err != nil {
		c.logger.Error(err)
		return helpers.ErrInternalError
	}

	return ctx.SendStatus(http.StatusOK)
}
