package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

type adminUsers struct {
	router   fiber.Router
	services *types.Services
}

func CreateAdminUsers(router fiber.Router, services *types.Services) fiber.Router {
	adminUsers := &adminUsers{
		services: services,
		router:   router,
	}

	return adminUsers.router.
		Group("users").
		Post("ignored", adminUsers.post)
}

type ignoredUserPostDto struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type ignoredUsersPostDto struct {
	Users []ignoredUserPostDto `json:"users"`
}

func (c *adminUsers) post(ctx *fiber.Ctx) error {
	dto := &ignoredUsersPostDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	err = c.postService(dto)
	if err == nil {
		return ctx.SendStatus(200)
	}

	return err
}
