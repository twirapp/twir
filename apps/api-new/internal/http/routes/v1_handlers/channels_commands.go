package v1_handlers

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/rand"
	"strconv"
)

func (c *Handlers) GetChannelsCommands(ctx *fiber.Ctx) error {
	r := rand.Intn(100)

	return ctx.SendString("v1 version commands for channel " + ctx.Params("channelId") + " random = " + strconv.Itoa(r))
}
