package donate_stream

import (
	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/events"
	"net/http"
)

func Setup(router fiber.Router, services types.Services) {
	router.Post("integrations/donate-stream/:integrationId", handlePost(services))
}

type postDto struct {
	Type     string `json:"type,omitempty"`
	Uid      string `json:"uid"`
	Message  string `json:"message"`
	Sum      string `json:"sum"`
	Nickname string `json:"nickname"`
}

func handlePost(services types.Services) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		redis := do.MustInvoke[*redis.Client](di.Provider)
		logger := do.MustInvoke[interfaces.Logger](di.Provider)

		integration := model.ChannelsIntegrations{}
		err := services.DB.Where("id = ?", ctx.Params("integrationId")).First(&integration).Error
		if err != nil {
			return fiber.NewError(http.StatusNotFound, "integration not found")
		}

		dto := &postDto{}
		err = middlewares.ValidateBody(
			ctx,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)

		if err != nil {
			return err
		}

		if dto.Type == "confirm" {
			value, err := redis.Get(ctx.Context(), "donate_stream_confirmation"+integration.ID).Result()
			if err != nil {
				logger.Error(err)
				return ctx.SendStatus(http.StatusInternalServerError)
			}

			return ctx.Status(http.StatusOK).SendString(value)
		}

		eventsGrpcClient := do.MustInvoke[events.EventsClient](di.Provider)

		eventsGrpcClient.Donate(ctx.Context(), &events.DonateMessage{
			BaseInfo: &events.BaseInfo{ChannelId: integration.ChannelID},
			UserName: dto.Nickname,
			Amount:   dto.Sum,
			Currency: "RUB",
			Message:  dto.Message,
		})

		return ctx.SendStatus(http.StatusOK)
	}
}
