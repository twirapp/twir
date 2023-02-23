package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/api/auth/user_defaults"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func createUserDefaults(ctx *fiber.Ctx) error {
	user := ctx.Locals("dbUser").(model.Users)

	go func() {
		user_defaults.CreateRoles(user.ID)
	}()

	return ctx.Next()
}
