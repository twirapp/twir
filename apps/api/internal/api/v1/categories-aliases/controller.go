package categories_aliases

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("categories-aliases")

	middleware.Get("/", get(services))
	middleware.Post("/", post(services))
	middleware.Delete(":categoryAliasId", delete(services))
	middleware.Patch(":categoryAliasId", patch(services))

	return middleware
}

// Game aliases godoc
// @Security ApiKeyAuth
// @Summary Get channel categories aliases list
// @Tags CategoriesAliases
// @Accept json
// @Produce json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success 200 {array} model.ChannelCategoryAlias
// @Failure 500 {object} types.DOCApiInternalError
// @Router /v1/channels/{channelId}/categories-aliases [get]
func get(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		categoriesAliases, err := handleGetCategoryAliases(c.Params("channelId"), services)
		if err != nil {
			return err
		}

		return c.JSON(categoriesAliases)
	}
}

// Category aliases godoc
// @Security ApiKeyAuth
// @Summary Create category alias
// @Tags CategoriesAliases
// @Accept json
// @Produce json
// @Param data body categoryAliasDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success 200 {object} model.ChannelCategoryAlias
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router /v1/channels/{channelId}/categories-aliases [post]
func post(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dto := &categoryAliasDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		categoryAlias, err := handlePost(c.Params("channelId"), dto, services)
		if err == nil && categoryAlias != nil {
			return c.JSON(categoryAlias)
		}

		return err
	}
}

// Category aliases godoc
// @Security ApiKeyAuth
// @Summary Delete category alias
// @Tags CategoriesAliases
// @Accept json
// @Produce json
// @Param        channelId   path      string  true  "ID of channel"
// @Param        categoryAliasId   path      string  true  "ID of category"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Router       /v1/channels/{channelId}/categories-aliases/{categoryAliasId} [delete]
func delete(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := handleDelete(c.Params("categoryAliasId"), services)
		if err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusOK)
	}
}

// Category aliases godoc
// @Security ApiKeyAuth
// @Summary Update category alias
// @Tags CategoriesAliases
// @Accept json
// @Produce json
// @Param data body categoryAliasDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Param        categoryAliasId   path      string  true  "ID of category"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/categories-aliases/{categoryAliasId} [patch]
func patch(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dto := &categoryAliasPatchDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		categoryAlias, err := handlePatch(c.Params("channelId"), c.Params("categoryAliasId"), dto, services)
		if err != nil {
			return err
		}

		return c.JSON(categoryAlias)
	}
}
