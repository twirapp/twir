package moderation

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("moderation")
	middleware.Get("", get(services))
	middleware.Post("", post(services))
	middleware.Post("/title", postTitle(services))
	middleware.Post("/category", postCategory(services))

	return middleware
}

// Moderation godoc
// @Security ApiKeyAuth
// @Summary      Get moderation settings
// @Tags         Moderation
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {array}  model.ChannelsModerationSettings
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/moderation [get]
func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		settings, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}

		return c.JSON(settings)
	}
}

// Moderation godoc
// @Security ApiKeyAuth
// @Summary      Create command
// @Tags         Moderation
// @Accept       json
// @Produce      json
// @Param data body moderationDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200  {array}  model.ChannelsModerationSettings
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/moderation [post]
func post(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := moderationDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			&dto,
		)
		if err != nil {
			return err
		}
		settings, err := handleUpdate(c.Params("channelId"), &dto, services)
		if err != nil {
			return err
		}

		return c.JSON(settings)
	}
}

func postTitle(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := postTitleDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			&dto,
		)
		if err != nil {
			return err
		}

		title, err := handlePostTitle(c.Params("channelId"), &dto, services)
		if err != nil {
			return err
		}

		return c.JSON(*title)
	}
}

func postCategory(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := postCategoryDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			&dto,
		)
		if err != nil {
			return err
		}

		category, err := handlePostCategory(c.Params("channelId"), &dto, services)
		if err != nil {
			return err
		}

		return c.JSON(category)
	}
}
