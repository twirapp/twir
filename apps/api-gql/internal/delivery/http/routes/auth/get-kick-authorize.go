package auth

import (
	"context"

	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

type kickAuthorizeResponse struct {
	AuthorizeURL string `json:"authorize_url"`
}

type kickAuthorizeOutput struct {
	Body kickAuthorizeResponse
}

type kickAuthorizeInput struct {
	RedirectTo string `query:"redirect_to"`
}

func (a *Auth) handleKickAuthorize(ctx context.Context, input kickAuthorizeInput) (*kickAuthorizeOutput, error) {
	return a.handlePlatformAuthorize(ctx, platformentity.PlatformKick, input.RedirectTo)
}

func (a *Auth) handlePlatformAuthorize(
	ctx context.Context,
	platform platformentity.Platform,
	redirectTo string,
) (*kickAuthorizeOutput, error) {
	authorizeURL, err := a.startPlatformAuth(ctx, platform, redirectTo)
	if err != nil {
		return nil, a.platformAuthHTTPError(err)
	}

	return &kickAuthorizeOutput{
		Body: kickAuthorizeResponse{
			AuthorizeURL: authorizeURL,
		},
	}, nil
}
