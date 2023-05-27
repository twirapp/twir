package auth_handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/tsuwari/apps/api-new/internal/http/helpers"
	"net/http"
)

func (c *AuthHandlers) GetLink(ctx *fiber.Ctx) error {
	twitchClient, err := helix.NewClient(&helix.Options{
		ClientID:    c.config.TwitchClientId,
		RedirectURI: c.config.TwitchCallbackUrl,
	})

	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	state := ctx.Query("state")
	if state == "" {
		return helpers.CreateBusinessErrorWithMessage(400, "state is missed")
	}

	url := twitchClient.GetAuthorizationURL(&helix.AuthorizationURLParams{
		ResponseType: "code",
		Scopes:       twitchScopes,
		State:        state,
		ForceVerify:  false,
	})

	return ctx.Redirect(url)
}
