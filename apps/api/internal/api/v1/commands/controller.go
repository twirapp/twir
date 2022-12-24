package commands

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
)

func Setup(router fiber.Router) fiber.Router {
	middleware := router.Group("commands")

	middleware.Get("", get)
	middleware.Post("", post)
	middleware.Delete(":commandId", delete)
	middleware.Put(":commandId", put)
	middleware.Patch(":commandId", patch)

	return router
}

type JSONResult struct{}

// Commands godoc
// @Security ApiKeyAuth
// @Summary      Get channel commands list
// @Tags         Commands
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId"
// @Success      200  {array}  model.ChannelsCommands
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/commands [get]
var get = func(c *fiber.Ctx) error {
	c.JSON(handleGet(c.Params("channelId")))

	return nil
}

// Commands godoc
// @Security ApiKeyAuth
// @Summary      Create command
// @Tags         Commands
// @Accept       json
// @Produce      json
// @Param data body commandDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200  {object}  model.ChannelsCommands
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/commands [post]
var post = func(c *fiber.Ctx) error {
	dto := &commandDto{}
	err := middlewares.ValidateBody(c, dto)
	if err != nil {
		return err
	}

	cmd, err := handlePost(c.Params("channelId"), dto)
	if err == nil {
		return c.JSON(cmd)
	}

	return err
}

// Commands godoc
// @Security ApiKeyAuth
// @Summary      Delete command
// @Tags         Commands
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ID of channel"
// @Param        commandId   path      string  true  "ID of command"
// @Success      200  {object}  model.ChannelsCommands
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/commands/{commandId} [delete]
var delete = func(c *fiber.Ctx) error {
	err := handleDelete(c.Params("channelId"), c.Params("commandId"))
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

// Commands godoc
// @Security ApiKeyAuth
// @Summary      Update command
// @Tags         Commands
// @Accept       json
// @Produce      json
// @Param data body commandDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Param        commandId   path      string  true  "ID of command"
// @Success      200  {object}  model.ChannelsCommands
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/commands/{commandId} [put]
var put = func(c *fiber.Ctx) error {
	dto := &commandDto{}
	err := middlewares.ValidateBody(c, dto)
	if err != nil {
		return err
	}

	cmd, err := handleUpdate(c.Params("channelId"), c.Params("commandId"), dto)
	if err == nil && cmd != nil {
		return c.JSON(cmd)
	}

	return err
}

// Commands godoc
// @Security ApiKeyAuth
// @Summary      Partialy update command
// @Tags         Commands
// @Accept       json
// @Produce      json
// @Param data body commandPatchDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Param        commandId   path      string  true  "ID of command"
// @Success      200  {object}  model.ChannelsCommands
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/commands/{commandId} [patch]
var patch = func(c *fiber.Ctx) error {
	dto := &commandPatchDto{}
	err := middlewares.ValidateBody(c, dto)
	if err != nil {
		return err
	}

	cmd, err := handlePatch(c.Params("channelId"), c.Params("commandId"), dto)
	if err == nil && cmd != nil {
		return c.JSON(cmd)
	}

	return err
}
