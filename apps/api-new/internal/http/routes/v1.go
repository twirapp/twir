package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-new/internal/http/routes/v1_handlers"
)

func NewV1(app *fiber.App, handlers *v1_handlers.Handlers) {
	group := app.Group("v1")

	group.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("v1 version root")
	})

	channelGroup := group.Group("channels/:channelId")

	commandsGroup := channelGroup.Group("commands")
	commandsGroup.Get("/", handlers.GetChannelsCommands)
	commandsGroup.Post("/", handlers.CreateCommand)
	commandsGroup.Delete("/:commandId", handlers.DeleteCommand)
}
