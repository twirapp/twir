package roles_users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group(":roleId/users")
	middleware.Get("/", getUsers())
	middleware.Put("/", updateUsers(services))

	return middleware
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
func getUsers() fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := getUsersService(c.Params("roleId"))
		if err != nil {
			return err
		}

		return c.JSON(users)
	}
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
func updateUsers(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dto := &roleUserDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		err = updateUsersService(c.Params("roleId"), dto.UserNames)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
