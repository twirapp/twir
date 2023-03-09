package commands

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/commands/commands_group"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

type Commands struct {
	services *types.Services
	router   fiber.Router
}

func NewController(router fiber.Router, services *types.Services) fiber.Router {
	commands := &Commands{
		services: services,
		router:   router,
	}

	commands.router = commands.router.
		Group("commands").
		Get("", commands.get).
		Post("", commands.post).
		Delete(":commandId", commands.delete).
		Put(":commandId", commands.put).
		Patch(":commandId", commands.patch)

	commands_group.NewCommandsGroup(commands.router, services)

	return commands.router
}

type JSONResult struct{}

// Commands godoc
// @Security ApiKeyAuth
// @Summary      Get channel commands list
// @Tags         Commands
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {array}  model.ChannelsCommands
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/commands [get]
func (c *Commands) get(ctx *fiber.Ctx) error {
	cmds, err := c.getService(ctx.Params("channelId"))
	if err != nil {
		return err
	}

	return ctx.JSON(cmds)
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
func (c *Commands) post(ctx *fiber.Ctx) error {
	dto := &commandDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	cmd, err := c.postService(ctx.Params("channelId"), dto)
	if err == nil {
		return ctx.JSON(cmd)
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
func (c *Commands) delete(ctx *fiber.Ctx) error {
	err := c.deleteService(ctx.Params("channelId"), ctx.Params("commandId"))
	if err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
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
func (c *Commands) put(ctx *fiber.Ctx) error {
	dto := &commandDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	cmd, err := c.putService(ctx.Params("channelId"), ctx.Params("commandId"), dto)
	if err != nil {
		return err
	}

	return ctx.JSON(cmd)
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
func (c *Commands) patch(ctx *fiber.Ctx) error {
	dto := &commandPatchDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	cmd, err := c.patchService(ctx.Params("channelId"), ctx.Params("commandId"), dto)
	if err == nil && cmd != nil {
		return ctx.JSON(cmd)
	}

	return err
}
