package donate_stream

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
	"net/http"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("donate-stream")
	middleware.Get("", get(services))
	middleware.Post(":integrationId", post(services))

	return middleware
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get Donate.stream id
// @Tags         Integrations|DonateStream
// @Accept       json
// @Produce      plain
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {string} string
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/donate-stream [get]
func get(services types.Services) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		data, err := handleGet(services, ctx.Params("channelId"))

		if err != nil {
			return err
		}

		return ctx.JSON(data)
	}
}

type postDto struct {
	Secret string `validate:"required" json:"secret"`
}

func post(services types.Services) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		dto := &postDto{}
		err := middlewares.ValidateBody(
			ctx,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		err = handlePost(services, ctx.Params("integrationId"), dto.Secret)
		if err != nil {
			return err
		}

		return ctx.SendStatus(http.StatusOK)
	}
}
