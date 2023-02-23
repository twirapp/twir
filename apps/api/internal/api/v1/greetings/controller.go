package greetings

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("greetings")

	middleware.Get("", get(services))
	middleware.Post("", post(services))
	middleware.Delete(":greetingId", delete(services))
	middleware.Put(":greetingId", put(services))
	middleware.Patch(":greetingId", patch(services))

	return middleware
}

// Greetings godoc
// @Security ApiKeyAuth
// @Summary      Get channel greeting list
// @Tags         Greetings
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {array}  Greeting
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/greetings [get]
func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(handleGet(c.Params("channelId"), services))
	}
}

// greetings godoc
// @Security ApiKeyAuth
// @Summary      Create greeting
// @Tags         Greetings
// @Accept       json
// @Produce      json
// @Param data body greetingsDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200  {object}  Greeting
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/greetings [post]
func post(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &greetingsDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		greeting, err := handlePost(c.Params("channelId"), dto, services)
		if err == nil && greeting != nil {
			return c.JSON(greeting)
		}

		return err
	}
}

// greetings godoc
// @Security ApiKeyAuth
// @Summary      Delete greeting
// @Tags         Greetings
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ID of channel"
// @Param        greetingId   path      string  true  "ID of greeting"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Router       /v1/channels/{channelId}/greetings/{greetingId} [delete]
func delete(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := handleDelete(c.Params("greetingId"), services)
		if err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusOK)
	}
}

// Greetings godoc
// @Security ApiKeyAuth
// @Summary      Update greeting
// @Tags         Greetings
// @Accept       json
// @Produce      json
// @Param data body greetingsDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Param        channelId   path      string  true  "ID of greeting"
// @Success      200  {object}  Greeting
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/greetings/{greetingId} [put]
func put(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &greetingsDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}
		greeting, err := handleUpdate(c.Params("greetingId"), dto, services)
		if err != nil {
			return err
		}

		return c.JSON(greeting)
	}
}

// Greetings godoc
// @Security ApiKeyAuth
// @Summary      Update greeting
// @Tags         Greetings
// @Accept       json
// @Produce      json
// @Param data body greetingsPatchDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Param        channelId   path      string  true  "ID of greeting"
// @Success      200  {object}  Greeting
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/greetings/{greetingId} [put]
func patch(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &greetingsPatchDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}
		greeting, err := handlePatch(c.Params("channelId"), c.Params("greetingId"), dto, services)
		if err != nil {
			return err
		}

		return c.JSON(greeting)
	}
}
