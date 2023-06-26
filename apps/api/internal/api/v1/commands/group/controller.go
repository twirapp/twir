package group

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/twir/apps/api/internal/middlewares"
	"github.com/satont/twir/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("groups")

	middleware.Get("/", getGroups(services))
	middleware.Post("/", createGroup(services))
	middleware.Delete(":groupId", deleteGroup(services))
	middleware.Put(":groupId", updateGroup(services))

	return middleware
}

// CommandsGroup godoc
// @Security ApiKeyAuth
// @Summary      Get channel commands groups list
// @Tags         CommandsGroup
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {array}  model.ChannelCommandGroup
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/commands/groups [get]
func getGroups(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		groups, err := getGroupsService(c.Params("channelId"), services)
		if err != nil {
			return err
		}

		return c.JSON(groups)
	}
}

// CommandsGroup
// @Security ApiKeyAuth
// @Summary      Create command group
// @Tags         CommandsGroup
// @Accept       json
// @Produce      json
// @Param data body groupDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200  {object}  model.ChannelCommandGroup
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/commands/groups [post]
func createGroup(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dto := &groupDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		group, err := createGroupService(c.Params("channelId"), dto)
		if err != nil {
			return err
		}

		return c.JSON(group)
	}
}

// CommandsGroup godoc
// @Security ApiKeyAuth
// @Summary      Delete command group
// @Tags         CommandsGroup
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ID of channel"
// @Param        groupId   path      string  true  "ID of group"
// @Success      204
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/commands/groups/{groupId} [delete]
func deleteGroup(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := deleteGroupService(c.Params("channelId"), c.Params("groupId"))
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// CommandsGroup godoc
// @Security ApiKeyAuth
// @Summary      Update command group
// @Tags         CommandsGroup
// @Accept       json
// @Produce      json
// @Param data body groupDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Param        groupId   path      string  true  "ID of group"
// @Success      200  {object}  model.ChannelCommandGroup
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/commands/groups/{groupId} [put]
func updateGroup(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dto := &groupDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		group, err := updateGroupService(c.Params("channelId"), c.Params("groupId"), dto)
		if err != nil {
			return err
		}

		return c.JSON(group)
	}
}
