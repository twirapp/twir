package donatepay

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("donatepay")
	middleware.Get("", get(services))
	middleware.Post("", post(services))

	return middleware
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get DonatePay token
// @Tags         Integrations|DonatePay
// @Accept       json
// @Produce      plain
// @Param        channelId   path      string  true  "ChannelId"
// @Success      200  {string} string
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/donatepay [get]
func get(services types.Services) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		data, err := handleGet(services, ctx.Params("channelId"))

		if err != nil {
			return err
		}

		return ctx.JSON(data)
	}
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Authorize DonatePay
// @Tags         Integrations|DonatePay
// @Accept       json
// @Param data body createOrUpdateDTO true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/donatepay [post]
func post(services types.Services) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		dto := &createOrUpdateDTO{}
		err := middlewares.ValidateBody(
			ctx,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		err = handlePost(services, ctx.Params("channelId"), dto)
		if err != nil {
			return err
		}

		return ctx.SendStatus(200)
	}
}
