package donatello

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
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
	return func(ctx *fiber.Ctx) error {
		apiKey := ctx.GetRespHeader("X-Key")

		if apiKey == "" {
			return ctx.SendStatus(http.StatusUnauthorized)
		}

		integration := model.ChannelsIntegrations{}
		err := services.DB.Where(`"apiKey" = ?`, apiKey).Find(&integration).Error
		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, "internal error")
		}

		if integration.ID == "" {
			return fiber.NewError(http.StatusNotFound, "not found")
		}

		dto := &postDto{}
		err = middlewares.ValidateBody(
			ctx,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)

		spew.Dump(dto)

		if err != nil {
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
