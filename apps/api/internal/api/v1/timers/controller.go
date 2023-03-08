package timers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services *types.Services) fiber.Router {
	middleware := router.Group("timers")
	middleware.Get("", get(services))
	middleware.Post("", post(services))
	middleware.Delete(":timerId", delete(services))
	middleware.Put(":timerId", put(services))
	middleware.Patch(":timerId", patch(services))

	return middleware
}

// Timers godoc
// @Security ApiKeyAuth
// @Summary      Get channel timers list
// @Tags         Timers
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {array}  model.ChannelsTimers
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/timers [get]
func get(services *types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		timers := handleGet(c.Params("channelId"), services)

		return c.JSON(timers)
	}
}

// Timers godoc
// @Security ApiKeyAuth
// @Summary      Create timer
// @Tags         Timers
// @Accept       json
// @Produce      json
// @Param data body timerDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200  {object}  model.ChannelsTimers
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/timers [post]
func post(services *types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &timerDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		cmd, err := handlePost(c.Params("channelId"), dto, services)
		if err == nil {
			return c.JSON(cmd)
		}

		return err
	}
}

// Timers godoc
// @Security ApiKeyAuth
// @Summary      Delete timer
// @Tags         Timers
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ID of channel"
// @Param        timerId   path      string  true  "ID of timer"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/timers/{timerId} [delete]
func delete(services *types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := handleDelete(c.Params("timerId"), services)
		if err != nil {
			return err
		}
		return c.SendStatus(200)
	}
}

// Timers godoc
// @Security ApiKeyAuth
// @Summary      Update timer
// @Tags         Timers
// @Accept       json
// @Produce      json
// @Param data body timerDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Param        timerId   path      string  true  "ID of timer"
// @Success      200  {object}  model.ChannelsTimers
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/timers/{timerId} [put]
func put(services *types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &timerDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		cmd, err := handlePut(c.Params("timerId"), dto, services)
		if err == nil {
			return c.JSON(cmd)
		}

		return err
	}
}

// Timers godoc
// @Security ApiKeyAuth
// @Summary      Partially update timer
// @Tags         Timers
// @Accept       json
// @Produce      json
// @Param data body timerPatchDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Param        timerId   path      string  true  "ID of timer"
// @Success      200  {object}  model.ChannelsTimers
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/timers/{timerId} [patch]
func patch(services *types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &timerPatchDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		cmd, err := handlePatch(c.Params("timerId"), dto, services)
		if err == nil {
			return c.JSON(cmd)
		}

		return err
	}
}
