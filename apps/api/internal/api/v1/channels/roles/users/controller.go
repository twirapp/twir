package roles_users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

type RolesUsers struct {
	services *types.Services
}

func NewRolesUsers(router fiber.Router, services *types.Services) fiber.Router {
	rolesUsers := &RolesUsers{
		services: services,
	}

	return router.Group(":roleId/users").
		Get("/", rolesUsers.get).
		Put("/", rolesUsers.put)
}

// Roles godoc
// @Security ApiKeyAuth
// @Summary Get users
// @Description Get users
// @Tags Roles
// @Accept json
// @Produce json
// @Param channelId path string true "Channel ID" default({{channelId}})
// @Param roleId path string true "Role ID" default({{roleId}})
// @Success      200  {array}  roleUser
// @Failure 404 {object} types.DOCApiNotFoundError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/roles/{roleId}/users [get]
func (c *RolesUsers) get(ctx *fiber.Ctx) error {
	users, err := c.getService(ctx.Params("roleId"))
	if err != nil {
		return err
	}

	return ctx.JSON(users)
}

// Roles godoc
// @Security ApiKeyAuth
// @Summary Update users
// @Description Update users
// @Tags Roles
// @Accept json
// @Produce json
// @Param channelId path string true "Channel ID" default({{channelId}})
// @Param roleId path string true "Role ID" default({{roleId}})
// @Param role body roleUserDto true "Role"
// @Success      200  {object}  model.ChannelRole
// @Failure 404 {object} types.DOCApiNotFoundError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/roles/{roleId}/users [put]
func (c *RolesUsers) put(ctx *fiber.Ctx) error {
	dto := &roleUserDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	err = c.putService(ctx.Params("roleId"), dto.UserNames)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
