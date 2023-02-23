package roles

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("/roles")
	middleware.Get("/", getRoles(services))
	middleware.Put("/:roleId", updateRole(services))
	middleware.Delete("/:roleId", deleteRole(services))
	middleware.Post("/", createRole(services))

	return middleware
}

// Roles godoc
// @Security ApiKeyAuth
// @Summary Get roles
// @Description Get roles
// @Tags Roles
// @Accept json
// @Produce json
// @Param channelId path string true "Channel ID" default({{channelId}})
// @Success      200  {array}  model.ChannelRole
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/roles [get]
func getRoles(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roles, err := getRolesService(c.Params("channelId"))
		if err != nil {
			return err
		}
		return c.JSON(roles)
	}
}

// Roles godoc
// @Security ApiKeyAuth
// @Summary Update role
// @Description Update role
// @Tags Roles
// @Accept json
// @Produce json
// @Param channelId path string true "Channel ID" default({{channelId}})
// @Param roleId path string true "Role ID" default({{roleId}})
// @Param role body roleDto true "Role"
// @Success      200  {object}  model.ChannelRole
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/roles/{roleId} [put]
func updateRole(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dto := &roleDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		newRole, err := updateRoleService(c.Params("channelId"), c.Params("roleId"), dto)
		if err != nil {
			return err
		}

		return c.JSON(newRole)
	}
}

// Roles godoc
// @Security ApiKeyAuth
// @Summary Delete role
// @Description Delete role
// @Tags Roles
// @Accept json
// @Produce json
// @Param channelId path string true "Channel ID" default({{channelId}})
// @Param roleId path string true "Role ID" default({{roleId}})
// @Success      204
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/roles/{roleId} [delete]
func deleteRole(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := deleteRoleService(c.Params("channelId"), c.Params("roleId"))
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// Roles godoc
// @Security ApiKeyAuth
// @Summary Create role
// @Description Create role
// @Tags Roles
// @Accept json
// @Produce json
// @Param channelId path string true "Channel ID" default({{channelId}})
// @Param role body roleDto true "Role"
// @Success      200  {object}  model.ChannelRole
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/roles [post]
func createRole(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dto := &roleDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		newRole, err := createRoleService(c.Params("channelId"), dto)
		if err != nil {
			return err
		}

		return c.JSON(newRole)
	}
}
