package community

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
	"net/http"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("community")
	middleware.Get("users", get(services))
	middleware.Delete("users/stats", resetStats(services))

	return middleware
}

// Community godoc
// @Security ApiKeyAuth
// @Summary      Get channel users list from database
// @Tags         Community
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {array}  User
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/community/users [get]
func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		limit := c.Query("limit", "50")
		page := c.Query("page", "1")
		sortBy := c.Query("sortBy", "watched")
		order := c.Query("order", "desc")

		users, err := handleGet(c.Params("channelId"), limit, page, sortBy, order)
		if err != nil {
			return err
		}

		return c.JSON(users)
	}
}

// Community godoc
// @Security ApiKeyAuth
// @Summary      Reset stats of channel
// @Tags         Community
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/community/users/stats [get]
func resetStats(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &resetStatsDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		err = handleResetStats(c.Params("channelId"), dto)
		if err != nil {
			return err
		}

		return c.SendStatus(http.StatusOK)
	}
}
