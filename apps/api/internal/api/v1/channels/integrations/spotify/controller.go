package spotify

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/satont/tsuwari/apps/api/internal/types"

	"github.com/satont/tsuwari/apps/api/internal/middlewares"
)

type Spotify struct {
	services *types.Services
}

func NewController(router fiber.Router, services *types.Services) fiber.Router {
	spotify := &Spotify{
		services,
	}

	return router.Group("spotify").
		Get(
			"", cache.New(cache.Config{
				Expiration: 31 * 24 * time.Hour,
				Storage:    services.RedisStorage,
				KeyGenerator: func(c *fiber.Ctx) string {
					return fmt.Sprintf("fiber:cache:integrations:spotify:profile:%s", c.Params("channelId"))
				},
			}),
			spotify.get,
		).
		Get("auth", spotify.authLink).
		Post("", spotify.post).
		Post("logout", spotify.logout)
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get Spotify profile
// @Tags         Integrations|Spotify
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {object}  spotify.SpotifyProfile
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/spotify [get]
func (c *Spotify) get(ctx *fiber.Ctx) error {
	profile, err := c.getService(ctx.Params("channelId"))
	if err != nil {
		return err
	}
	return ctx.JSON(profile)
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get Spotify auth link
// @Tags         Integrations|Spotify
// @Accept       json
// @Produce      plain
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success 200 {string} string	"Auth link"
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/spotify/auth [get]
func (c *Spotify) authLink(ctx *fiber.Ctx) error {
	authLink, err := c.getAuthLinkService()
	if err != nil {
		return err
	}

	return ctx.SendString(*authLink)
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Authorize Spotify
// @Tags         Integrations|Spotify
// @Accept       json
// @Produce      json
// @Param data body tokenDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/spotify [post]
func (c *Spotify) post(ctx *fiber.Ctx) error {
	dto := &tokenDto{}
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
		fmt.Sprintf("fiber:cache:integrations:spotify:profile:%s", channelId),
		"GET",
	)

	return ctx.SendStatus(200)
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Logout
// @Tags         Integrations|Spotify
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404 {object} types.DOCApiBadRequest
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/spotify/logout [post]
func (c *Spotify) logout(ctx *fiber.Ctx) error {
	channelId := ctx.Params("channelId")
	err := c.logoutService(channelId)
	if err != nil {
		return err
	}

	c.services.RedisStorage.DeleteByMethod(
		fmt.Sprintf("fiber:cache:integrations:spotify:profile:%s", channelId),
		"GET",
	)

	return ctx.SendStatus(200)

}
