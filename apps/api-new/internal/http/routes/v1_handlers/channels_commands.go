package v1_handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/rand"
	"strconv"
)

func (c *Handlers) GetChannelsCommands(ctx *fiber.Ctx) error {
	r := rand.Intn(100)

	session, err := c.sessionStorage.Get(ctx)
	if err != nil {
		return err
	}

	fmt.Println(session.Get("random"))
	if session.Get("random") == nil {
		session.Set("random", r)
		session.Save()
	}

	return ctx.SendString("v1 version commands for channel " + ctx.Params("channelId") + " random = " + strconv.Itoa(session.Get("random").(int)))
}
