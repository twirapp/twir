package donatello

import (
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
	router.Post("integrations/donatello", handlePost(services))
}

type postDto struct {
	PubId       string `json:"pubId"`
	ClientName  string `json:"clientName"`
	Message     string `json:"message"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Source      string `json:"source"`
	Goal        string `json:"goal"`
	IsPublished bool   `json:"isPublished"`
	CreatedAt   string `json:"createdAt"`
}

func handlePost(services types.Services) fiber.Handler {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	return func(ctx *fiber.Ctx) error {
		apiKey := ctx.Get("X-Key")

		if apiKey == "" {
			logger.Infow("No key", "key", apiKey)
			return ctx.SendStatus(http.StatusUnauthorized)
		}

		integration := &model.ChannelsIntegrations{}
		err := services.DB.Where(`"id" = ?`, apiKey).Find(integration).Error
		if err != nil {
			logger.Error(err)
			return fiber.NewError(http.StatusInternalServerError, "internal error")
		}

		if integration.ID == "" {
			logger.Infow("No integration", "key", apiKey)
			return fiber.NewError(http.StatusNotFound, "not found")
		}

		dto := &postDto{}
		err = middlewares.ValidateBody(
			ctx,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)

		if err != nil {
			logger.Error(err)
			return err
		}

		eventsGrpcClient := do.MustInvoke[events.EventsClient](di.Provider)

		eventsGrpcClient.Donate(ctx.Context(), &events.DonateMessage{
			BaseInfo: &events.BaseInfo{ChannelId: integration.ChannelID},
			UserName: dto.ClientName,
			Amount:   dto.Amount,
			Currency: dto.Currency,
			Message:  dto.Message,
		})

		return ctx.SendStatus(http.StatusOK)
	}
}
