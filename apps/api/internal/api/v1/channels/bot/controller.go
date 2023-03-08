package bot

import (
	"fmt"
	"time"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"

	_ "github.com/satont/tsuwari/libs/types/types/api/bot"
)

type Bot struct {
	services *types.Services
	router   fiber.Router
}

func NewBot(router fiber.Router, services *types.Services) fiber.Router {
	bot := &Bot{
		services: services,
		router:   router,
	}

	return bot.router.
		Group("bot").
		Get(
			"",
			cache.New(cache.Config{
				Expiration: 20 * time.Second,
				Storage:    services.RedisStorage,
				KeyGenerator: func(c *fiber.Ctx) string {
					return fmt.Sprintf("channels:isBotMod:%s", c.Params("channelId"))
				},
			}),
			bot.get,
		).
		Patch(
			"connection",
			limiter.New(limiter.Config{
				Max:        2,
				Expiration: 1 * time.Second,
				KeyGenerator: func(c *fiber.Ctx) string {
					dbUser := c.Locals("dbUser").(model.Users)
					return fmt.Sprintf("fiber:limiter:bot:connection:%s", dbUser.ID)
				},
				LimitReached: func(c *fiber.Ctx) error {
					header := c.GetRespHeader("Retry-After", "2")
					return c.Status(429).JSON(fiber.Map{"message": fmt.Sprintf("wait %s seconds", header)})
				},
				Storage: services.RedisStorage,
			}),
			bot.patch,
		)
}

// Bot godoc
// @Security ApiKeyAuth
// @Summary      Check does bot moderator on the channel
// @Tags         Bot
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {object}  BotInfo
// @Failure 404 {object} types.DOCApiInternalError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/bot/checkmod [get]
func (c *Bot) get(ctx *fiber.Ctx) error {
	isBotMod, err := c.getService(ctx.Params("channelId"))
	if err != nil {
		return err
	}

	if isBotMod == nil {
		return ctx.JSON(false)
	}

	return ctx.JSON(*isBotMod)
}

// Bot godoc
// @Security ApiKeyAuth
// @Summary      Check does bot moderator on the channel
// @Tags         Bot
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Param data body connectionDto true "Data"
// @Success      200  {boolean}  boolean
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/bot/connection [patch]
func (c *Bot) patch(ctx *fiber.Ctx) error {
	dto := &connectionDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	channelId := ctx.Params("channelId")
	err = c.patchService(channelId, dto)

	if err != nil {
		return err
	}

	c.services.RedisStorage.DeleteByMethod(
		fmt.Sprintf("channels:isBotMod:%s", channelId),
		"GET",
	)

	return ctx.SendStatus(200)
}
