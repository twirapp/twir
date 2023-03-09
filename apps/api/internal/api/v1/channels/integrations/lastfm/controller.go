package lastfm

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

type LastFM struct {
	services *types.Services
}

func NewController(router fiber.Router, services *types.Services) fiber.Router {
	lastFm := &LastFM{
		services,
	}

	return router.Group("lastfm").
		Get("auth", lastFm.auth).
		Post("", lastFm.post).
		Post("logout", lastFm.logout).
		Get(
			"",
			cache.New(cache.Config{
				Expiration: 31 * 24 * time.Hour,
				Storage:    services.RedisStorage,
				KeyGenerator: func(c *fiber.Ctx) string {
					return fmt.Sprintf("fiber:cache:integrations:lastfm:profile:%s", c.Params("channelId"))
				},
			}),
			lastFm.get,
		)

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
func (c *LastFM) get(ctx *fiber.Ctx) error {
	integration, err := c.getService(ctx.Params("channelId"))
	if err != nil {
		return err
	}
	return ctx.JSON(integration)
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
func (c *LastFM) auth(ctx *fiber.Ctx) error {
	integration, err := c.authService()
	if err != nil {
		return err
	}
	return ctx.JSON(integration)
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
func (c *LastFM) post(ctx *fiber.Ctx) error {
	dto := &lastfmDto{}
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
	err = c.postService(channelId, dto)

	if err != nil {
		return err
	}

	c.services.RedisStorage.DeleteByMethod(
		fmt.Sprintf("fiber:cache:integrations:lastfm:profile:%s", channelId),
		"GET",
	)

	return ctx.SendStatus(200)
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
func (c *LastFM) logout(ctx *fiber.Ctx) error {
	channelId := ctx.Params("channelId")
	err := c.logoutService(channelId)
	if err != nil {
		return err
	}

	c.services.RedisStorage.DeleteByMethod(
		fmt.Sprintf("fiber:cache:integrations:lastfm:profile:%s", channelId),
		"GET",
	)

	return ctx.SendStatus(200)
}
