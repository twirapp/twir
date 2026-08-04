package mcp_oauth

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"io"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	authsessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
	config "github.com/twirapp/twir/libs/config"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
)

type oauthService interface {
	RegisterClient(context.Context, service.RegisterClientInput) (entity.Client, error)
	GetClient(context.Context, string) (entity.Client, error)
	ValidateAuthorizeInput(context.Context, service.AuthorizeInput) (service.AuthorizationRequest, error)
	CreateAuthorizationCode(context.Context, service.CreateAuthorizationCodeInput) (service.IssuedAuthorizationCode, error)
	ExchangeAuthorizationCode(context.Context, service.ExchangeAuthorizationCodeInput) (service.TokenSet, error)
	Refresh(context.Context, service.RefreshInput) (service.TokenSet, error)
	Revoke(context.Context, service.RevokeInput) error
}

type sessionStore interface {
	SetMCPOAuthAttempt(context.Context, string, authsessions.MCPOAuthAttempt) error
	GetMCPOAuthAttempt(context.Context, string) (authsessions.MCPOAuthAttempt, error)
	DeleteMCPOAuthAttempt(context.Context, string) error
	GetInternalUserID(context.Context) (uuid.UUID, error)
	GetSelectedDashboard(context.Context) (string, error)
}

type registerRateLimiter interface {
	Use(context.Context, *rate_limiter.LeakyOptions, int) (*rate_limiter.LeakyResponse, error)
}

type Dependencies struct {
	Service             oauthService
	Sessions            sessionStore
	SiteBaseURL         string
	Random              io.Reader
	RegisterRateLimiter registerRateLimiter
}

type Handler struct {
	service             oauthService
	sessions            sessionStore
	origin              string
	random              io.Reader
	registerRateLimiter registerRateLimiter
}

func New(deps Dependencies) (*Handler, error) {
	if deps.Service == nil || deps.Sessions == nil {
		return nil, fmt.Errorf("MCP OAuth handler dependencies are required")
	}
	origin, err := canonicalOrigin(deps.SiteBaseURL)
	if err != nil {
		return nil, err
	}
	if deps.Random == nil {
		deps.Random = cryptorand.Reader
	}
	return &Handler{service: deps.Service, sessions: deps.Sessions, origin: origin, random: deps.Random, registerRateLimiter: deps.RegisterRateLimiter}, nil
}

func NewFromOpts(oauthService *service.Service, sessions *authsessions.Auth, config config.Config, rateLimiter *rate_limiter.LeakyBucketRateLimiter) (*Handler, error) {
	return New(Dependencies{Service: oauthService, Sessions: sessions, SiteBaseURL: config.SiteBaseUrl, RegisterRateLimiter: rateLimiter})
}

type route interface {
	GetMeta() huma.Operation
	Register(huma.API)
}

func (handler *Handler) routes() []route {
	return []route{
		newProtectedResourceMetadata(handler),
		newAuthorizationServerMetadata(handler),
		newScopeCatalog(),
		newRegisterOptions(handler),
		newRegisterClient(handler),
		newAuthorize(handler),
		newGetConsent(handler),
		newPostConsent(handler),
		newTokenOptions(handler),
		newToken(handler),
		newRevokeOptions(handler),
		newRevoke(handler),
	}
}

func (handler *Handler) Register(api huma.API) {
	for _, route := range handler.routes() {
		route.Register(api)
	}
}

type Registration struct{}

func RegisterRoutes(api huma.API, handler *Handler) Registration {
	handler.Register(api)
	return Registration{}
}

type emptyInput struct{}
type preflightOutput struct{ Status int }

func canonicalOrigin(raw string) (string, error) {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("invalid MCP OAuth site base URL")
	}
	return (&url.URL{Scheme: parsed.Scheme, Host: parsed.Host}).String(), nil
}
