package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-new/internal/http/routes/auth_handlers"
)

func NewAuth(app *fiber.App, handlers *auth_handlers.AuthHandlers) {
	group := app.Group("auth")

	group.Get("/", handlers.GetLink)
	group.Post("/", handlers.PostCode)
	group.Delete("/", handlers.Logout)
	group.Get("/profile", handlers.GetProfile)
}
