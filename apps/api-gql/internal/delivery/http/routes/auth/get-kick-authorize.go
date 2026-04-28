package auth

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	kickplatform "github.com/twirapp/twir/apps/api-gql/internal/platform/kick"
)

const kickCodeVerifierSessionKey = "kick_code_verifier"

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
	if a.kickProvider == nil {
		return nil, huma.Error500InternalServerError("Kick provider is not configured", fmt.Errorf("kick provider is nil"))
	}

	codeVerifier, err := kickplatform.GenerateCodeVerifier()
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot generate code verifier", err)
	}

	codeChallenge := kickplatform.GenerateCodeChallenge(codeVerifier)

	redirectTo := input.RedirectTo
	if redirectTo == "" {
		redirectTo = "/dashboard"
	}
	state := base64.StdEncoding.EncodeToString([]byte(redirectTo))

	a.sessions.Put(ctx, kickCodeVerifierSessionKey, codeVerifier)

	if err := a.sessions.Commit(ctx); err != nil {
		return nil, huma.Error500InternalServerError("Cannot commit session", err)
	}

	return &kickAuthorizeOutput{
		Body: kickAuthorizeResponse{
			AuthorizeURL: a.kickProvider.GetAuthURL(state, codeChallenge),
		},
	}, nil
}
