package tts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
	"net/http"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("tts")
	middleware.Use(middlewares.CheckUserAuth(services))
	middleware.Get("say", getSpeak(services))

	return middleware
}

func getSpeak(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		volume := c.Query("volume", "100")
		pitch := c.Query("pitch", "50")
		rate := c.Query("rate", "50")
		text := c.Query("text")
		voice := c.Query("voice")

		if voice == "" || text == "" {
			return c.SendStatus(http.StatusBadRequest)
		}

		r, err := handleGetSay(voice, pitch, volume, rate, text)
		if err != nil {
			return err
		}

		return c.Type("audio/ogg").SendStream(r)
	}
}
