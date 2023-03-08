package keywords

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

type Keywords struct {
	services *types.Services
}

func NewKeywords(router fiber.Router, services *types.Services) fiber.Router {
	keywords := &Keywords{
		services: services,
	}

	return router.Group("keywords").
		Get("", keywords.get).
		Post("", keywords.post).
		Delete(":keywordId", keywords.delete).
		Put(":keywordId", keywords.put).
		Patch(":keywordId", keywords.patch)
}

// Keywords godoc
// @Security ApiKeyAuth
// @Summary      Get channel keywords list
// @Tags         Keywords
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {array}  model.ChannelsKeywords
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/keywords [get]
func (c *Keywords) get(ctx *fiber.Ctx) error {
	keywords, err := c.getService(ctx.Params("channelId"))
	if err != nil {
		return err
	}
	return ctx.JSON(keywords)
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
func (c *Keywords) post(ctx *fiber.Ctx) error {
	dto := &keywordDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	keyword, err := c.postService(ctx.Params("channelId"), dto)
	if err == nil {
		return ctx.JSON(keyword)
	}

	return err
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
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/keywords/{keywordId} [delete]
func (c *Keywords) delete(ctx *fiber.Ctx) error {
	err := c.deleteService(ctx.Params("keywordId"))
	if err != nil {
		return err
	}

	return ctx.SendStatus(200)
}

// Keywords godoc
// @Security ApiKeyAuth
// @Summary      Update keyword
// @Tags         Keywords
// @Accept       json
// @Produce      json
// @Param data body keywordDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Param        keywordId   path      string  true  "ID of keyword"
// @Success      200  {object}  model.ChannelsKeywords
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/keywords/{keywordId} [put]
func (c *Keywords) put(ctx *fiber.Ctx) error {
	dto := &keywordDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	keyword, err := c.putService(ctx.Params("keywordId"), dto)
	if err != nil {
		return err
	}

	return ctx.JSON(keyword)
}

// Keywords godoc
// @Security ApiKeyAuth
// @Summary      Partially update keyword
// @Tags         Commands
// @Accept       json
// @Produce      json
// @Param data body keywordPatchDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Param        keywordId   path      string  true  "ID of keyword"
// @Success      200  {object}  model.ChannelsKeywords
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/keywords/{keywordId} [patch]
func (c *Keywords) patch(ctx *fiber.Ctx) error {
	dto := &keywordPatchDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	cmd, err := c.patchService(ctx.Params("channelId"), ctx.Params("keywordId"), dto)
	if err == nil && cmd != nil {
		return ctx.JSON(cmd)
	}

	return err
}
