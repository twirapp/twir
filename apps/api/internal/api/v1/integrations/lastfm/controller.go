package lastfm

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services *types.Services) fiber.Router {
	middleware := router.Group("lastfm")
	middleware.Get("auth", auth(services))
	middleware.Post("", post(services))
	middleware.Post("logout", logout(services))

	profileCache := cache.New(cache.Config{
		Expiration: 31 * 24 * time.Hour,
		Storage:    services.RedisStorage,
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf("fiber:cache:integrations:lastfm:profile:%s", c.Params("channelId"))
		},
	})

	middleware.Get("", profileCache, get(services))

	return middleware
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get LastFm profile
// @Tags         Integrations|Lastfm
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {array}  LastfmProfile
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/lastfm [get]
func get(services *types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		integration, err := handleProfile(c.Params("channelId"), services)
		if err != nil {
			return err
		}
		return c.JSON(integration)
	}
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get LastFm auth link
// @Tags         Integrations|Lastfm
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {string}  string
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/lastfm/auth [get]
func auth(services *types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		integration, err := handleAuth(services)
		if err != nil {
			return err
		}
		return c.JSON(integration)
	}
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Update LastFM
// @Tags         Integrations|Lastfm
// @Accept       json
// @Produce      json
// @Param data body lastfmDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200  {object} model.ChannelsIntegrations
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/lastfm [post]
func post(services *types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &lastfmDto{}
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
		err = handlePost(channelId, dto, services)

		if err != nil {
			return err
		}

		services.RedisStorage.DeleteByMethod(
			fmt.Sprintf("fiber:cache:integrations:lastfm:profile:%s", channelId),
			"GET",
		)

		return c.SendStatus(200)
	}
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Logout
// @Tags         Integrations|Lastfm
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404 {object} types.DOCApiBadRequest
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/lastfm/logout [post]
func logout(services *types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		channelId := c.Params("channelId")
		err := handleLogout(channelId, services)
		if err != nil {
			return err
		}

		services.RedisStorage.DeleteByMethod(
			fmt.Sprintf("fiber:cache:integrations:lastfm:profile:%s", channelId),
			"GET",
		)

		return c.SendStatus(200)
	}
}
