package variables

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("variables")
	middleware.Get("", get(services))
	middleware.Post("", post(services))
	middleware.Delete(":variableId", delete(services))
	middleware.Put(":variableId", put(services))
	middleware.Get("builtin", getBuiltIn(services))

	return middleware
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
func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		variables, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}
		return c.JSON(variables)
	}
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
func getBuiltIn(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		variables, err := handleGetBuiltIn(services)
		if err != nil {
			return err
		}

		return c.JSON(variables)
	}
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
func post(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &variableDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		variable, err := handlePost(c.Params("channelId"), dto, services)
		if err == nil {
			return c.JSON(variable)
		}

		return err
	}
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
func delete(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := handleDelete(c.Params("channelId"), c.Params("variableId"), services)
		if err != nil {
			return err
		}
		return c.SendStatus(200)
	}
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
func put(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &variableDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		variable, err := handleUpdate(c.Params("channelId"), c.Params("variableId"), dto, services)
		if err == nil {
			return c.JSON(variable)
		}

		return err
	}
}
