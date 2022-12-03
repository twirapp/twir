package keywords

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("keywords")

	middleware.Get("", get(services))
	middleware.Post("", post(services))
	middleware.Delete(":keywordId", delete(services))
	middleware.Put(":keywordId", put(services))

	return middleware
}

// Keywords godoc
// @Security ApiKeyAuth
// @Summary      Get channel keywords list
// @Tags         Keywords
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId"
// @Success      200  {array}  model.ChannelsKeywords
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/keywords [get]
func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		keywords, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}
		return c.JSON(keywords)
	}
}

// Keywords godoc
// @Security ApiKeyAuth
// @Summary      Create keyword
// @Tags         Keywords
// @Accept       json
// @Produce      json
// @Param data body keywordDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200  {object}  model.ChannelsKeywords
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/keywords [post]
func post(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &keywordDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		keyword, err := handlePost(c.Params("channelId"), dto, services)
		if err == nil {
			return c.JSON(keyword)
		}

		return err
	}
}

// Keywords godoc
// @Security ApiKeyAuth
// @Summary      Delete keyword
// @Tags         Keywords
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ID of channel"
// @Param        keywordId   path      string  true  "ID of keyword"
// @Success      200  {object}  model.ChannelsKeywords
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/keywords/{keywordId} [delete]
func delete(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := handleDelete(c.Params("keywordId"), services)
		if err != nil {
			return err
		}

		return c.SendStatus(200)
	}
}

// Keywords godoc
// @Security ApiKeyAuth
// @Summary      Update command
// @Tags         Keywords
// @Accept       json
// @Produce      json
// @Param data body keywordDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Param        keywordId   path      string  true  "ID of keyword"
// @Success      200  {object}  model.ChannelsKeywords
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Failute 404
// @Router       /v1/channels/{channelId}/keywords/{keywordId} [put]
func put(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &keywordDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		keyword, err := handleUpdate(c.Params("keywordId"), dto, services)
		if err != nil {
			return err
		}

		return c.JSON(keyword)
	}
}
