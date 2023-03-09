package variables

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

type Variables struct {
	services *types.Services
}

func NewController(router fiber.Router, services *types.Services) fiber.Router {
	variables := &Variables{}

	return router.Group("variables").
		Get("", variables.get).
		Post("", variables.post).
		Delete(":variableId", variables.delete).
		Put(":variableId", variables.put).
		Get("builtin", variables.builtIn)
}

// Variables godoc
// @Security ApiKeyAuth
// @Summary      Get channel variables list
// @Tags         Variables
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {array}  model.ChannelsCustomvars
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/variables [get]
func (c *Variables) get(ctx *fiber.Ctx) error {
	variables, err := c.getService(ctx.Params("channelId"))
	if err != nil {
		return err
	}
	return ctx.JSON(variables)
}

// Variables godoc
// @Security ApiKeyAuth
// @Summary      Get built-in variables
// @Tags         Variables
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {array}  parser.GetVariablesResponse_Variable
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/variables/builtin [get]
func (c *Variables) builtIn(ctx *fiber.Ctx) error {
	variables, err := c.builtInService()
	if err != nil {
		return err
	}

	return ctx.JSON(variables)
}

// Variables godoc
// @Security ApiKeyAuth
// @Summary      Create variable
// @Tags         Variables
// @Accept       json
// @Produce      json
// @Param data body variableDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200  {object}  model.ChannelsCustomvars
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/variables [post]
func (c *Variables) post(ctx *fiber.Ctx) error {
	dto := &variableDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	variable, err := c.postService(ctx.Params("channelId"), dto)
	if err == nil {
		return ctx.JSON(variable)
	}

	return err
}

// Variables godoc
// @Security ApiKeyAuth
// @Summary      Delete variable
// @Tags         Variables
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ID of channel"
// @Param        variableId   path      string  true  "ID of variable"
// @Success      200  {object}  model.ChannelsCustomvars
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/variables/{variableId} [delete]
func (c *Variables) delete(ctx *fiber.Ctx) error {
	err := c.deleteService(ctx.Params("channelId"), ctx.Params("variableId"))
	if err != nil {
		return err
	}
	return ctx.SendStatus(200)
}

// Variables godoc
// @Security ApiKeyAuth
// @Summary      Update variable
// @Tags         Variables
// @Accept       json
// @Produce      json
// @Param data body variableDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Param        variableId   path      string  true  "ID of variable"
// @Success      200  {object}  model.ChannelsCustomvars
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/variables/{variableId} [put]
func (c *Variables) put(ctx *fiber.Ctx) error {
	dto := &variableDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	variable, err := c.putService(ctx.Params("channelId"), ctx.Params("variableId"), dto)
	if err == nil {
		return ctx.JSON(variable)
	}

	return err
}
