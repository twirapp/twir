package auth

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	kickplatform "github.com/twirapp/twir/apps/api-gql/internal/platform/kick"
)

const kickCodeVerifierSessionKey = "kick_code_verifier"

type kickAuthorizeResponse struct {
	AuthorizeURL string `json:"authorize_url"`
}

type kickAuthorizeOutput struct {
	Body kickAuthorizeResponse
}

func (a *Auth) handleKickAuthorize(ctx context.Context) (*kickAuthorizeOutput, error) {
	codeVerifier, err := kickplatform.GenerateCodeVerifier()
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot generate code verifier", err)
	}

	codeChallenge := kickplatform.GenerateCodeChallenge(codeVerifier)
	state := uuid.NewString()

	a.sessions.Put(ctx, kickCodeVerifierSessionKey, codeVerifier)

	if a.kickProvider == nil {
		return nil, huma.Error500InternalServerError("Kick provider is not configured", fmt.Errorf("kick provider is nil"))
	}

	return &kickAuthorizeOutput{
		Body: kickAuthorizeResponse{
			AuthorizeURL: a.kickProvider.GetAuthURL(state, codeChallenge),
		},
	}, nil
}
