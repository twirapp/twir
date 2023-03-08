package roles

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/roles/users"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

type Roles struct {
	services *types.Services
}

func NewRoles(router fiber.Router, services *types.Services) fiber.Router {
	roles := Roles{services: services}

	middleware := router.Group("/roles")
	roles_users.NewRolesUsers(middleware, services)

	middleware.Get("/", roles.get).
		Put("/:roleId", roles.put).
		Delete("/:roleId", roles.delete).
		Post("/", roles.post)

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
func (c *Roles) get(ctx *fiber.Ctx) error {
	roles, err := c.getService(ctx.Params("channelId"))
	if err != nil {
		return err
	}
	return ctx.JSON(roles)
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
// @Failure 404 {object} types.DOCApiNotFoundError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/roles/{roleId} [put]
func (c *Roles) put(ctx *fiber.Ctx) error {
	dto := &roleDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	newRole, err := c.putService(ctx.Params("channelId"), ctx.Params("roleId"), dto)
	if err != nil {
		return err
	}

	return ctx.JSON(newRole)
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
// @Failure 404 {object} types.DOCApiNotFoundError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/roles/{roleId} [delete]
func (c *Roles) delete(ctx *fiber.Ctx) error {
	err := c.deleteService(ctx.Params("channelId"), ctx.Params("roleId"))
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
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
func (c *Roles) post(ctx *fiber.Ctx) error {
	dto := &roleDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	newRole, err := c.postService(ctx.Params("channelId"), dto)
	if err != nil {
		return err
	}

	return ctx.JSON(newRole)
}
