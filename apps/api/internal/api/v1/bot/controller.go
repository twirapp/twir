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

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("bot")

	botInfoCache := cache.New(cache.Config{
		Expiration: 20 * time.Second,
		Storage:    services.RedisStorage,
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf("channels:isBotMod:%s", c.Params("channelId"))
		},
	})

	middleware.Get("", botInfoCache, get(services))

	limit := limiter.New(limiter.Config{
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
	})

	middleware.Patch("connection", limit, patch(services))

	return middleware
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
func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		isBotMod, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}

		if isBotMod == nil {
			return c.JSON(false)
		}

		return c.JSON(*isBotMod)
	}
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
func patch(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &connectionDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		channelId := c.Params("channelId")
		err = handlePatch(channelId, dto, services)

		if err != nil {
			return err
		}

		services.RedisStorage.DeleteByMethod(
			fmt.Sprintf("channels:isBotMod:%s", channelId),
			"GET",
		)

		return c.SendStatus(200)
	}
}
