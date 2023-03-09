package donationalerts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

type DonationAlerts struct {
	services *types.Services
}

func NewController(router fiber.Router, services *types.Services) fiber.Router {
	donationAlerts := &DonationAlerts{
		services,
	}

	return router.Group("donationalerts").
		Get("auth", donationAlerts.getAuth).
		Get("", donationAlerts.get).
		Post("logout", donationAlerts.logout).
		Post("", donationAlerts.post)
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get DonationAlerts integration
// @Tags         Integrations|DonationAlerts
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {object}  model.ChannelsIntegrationsData
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/donationalerts [get]
func (c *DonationAlerts) get(ctx *fiber.Ctx) error {
	integration, err := c.getService(ctx.Params("channelId"))
	if err != nil {
		return err
	}
	return ctx.JSON(integration)
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get DonationAlerts auth link
// @Tags         Integrations|DonationAlerts
// @Accept       json
// @Produce      plain
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success 200 {string} string	"Auth link"
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/donationalerts/auth [get]
func (c *DonationAlerts) getAuth(ctx *fiber.Ctx) error {
	authLink, err := c.getAuthService()
	if err != nil {
		return err
	}

	return ctx.SendString(*authLink)
}

type tokenDto struct {
	Code string `validate:"required" json:"code"`
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Update auth of DonationAlerts
// @Tags         Integrations|DonationAlerts
// @Accept       json
// @Produce      json
// @Param data body tokenDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/donationalerts [post]
func (c *DonationAlerts) post(ctx *fiber.Ctx) error {
	dto := &tokenDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	err = c.postService(ctx.Params("channelId"), dto)
	if err != nil {
		return err
	}

	return ctx.SendStatus(200)
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Logout
// @Tags         Integrations|DonationAlerts
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404 {object} types.DOCApiBadRequest
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/donationalerts/logout [post]
func (c *DonationAlerts) logout(ctx *fiber.Ctx) error {
	channelId := ctx.Params("channelId")
	err := c.logoutService(channelId)
	if err != nil {
		return err
	}

	return ctx.SendStatus(200)
}
