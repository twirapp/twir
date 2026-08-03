package mcp_oauth

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	authsessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/server"
	servermiddlewares "github.com/twirapp/twir/apps/api-gql/internal/server/middlewares"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
	config "github.com/twirapp/twir/libs/config"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
	"go.uber.org/fx"
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

type Dependencies struct {
	Service           oauthService
	Sessions          sessionStore
	SiteBaseURL       string
	Random            io.Reader
	RegisterRateLimit gin.HandlerFunc
}

type Handler struct {
	service           oauthService
	sessions          sessionStore
	origin            string
	random            io.Reader
	registerRateLimit gin.HandlerFunc
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
	if deps.RegisterRateLimit == nil {
		deps.RegisterRateLimit = func(context *gin.Context) { context.Next() }
	}
	return &Handler{service: deps.Service, sessions: deps.Sessions, origin: origin, random: deps.Random, registerRateLimit: deps.RegisterRateLimit}, nil
}

type FxOpts struct {
	fx.In
	Service     *service.Service
	Sessions    *authsessions.Auth
	Config      config.Config
	Middlewares *servermiddlewares.Middlewares
}

func NewFx(opts FxOpts) (*Handler, error) {
	return New(Dependencies{Service: opts.Service, Sessions: opts.Sessions, SiteBaseURL: opts.Config.SiteBaseUrl, RegisterRateLimit: opts.Middlewares.RateLimit("mcp-oauth-register", 20, time.Minute)})
}

func (handler *Handler) Register(router gin.IRouter) {
	router.GET("/.well-known/oauth-protected-resource", handler.protectedResourceMetadata)
	router.GET("/.well-known/oauth-authorization-server", handler.authorizationServerMetadata)
	router.OPTIONS("/oauth/register", handler.publicPreflight)
	router.POST("/oauth/register", handler.publicCORS, handler.registerRateLimit, handler.register)
	router.GET("/oauth/authorize", handler.authorize)
	router.GET("/oauth/consent", handler.getConsent)
	router.POST("/oauth/consent", handler.postConsent)
	router.OPTIONS("/oauth/token", handler.publicPreflight)
	router.POST("/oauth/token", handler.publicCORS, handler.token)
	router.OPTIONS("/oauth/revoke", handler.publicPreflight)
	router.POST("/oauth/revoke", handler.publicCORS, handler.revoke)
}

func Register(router *server.Server, handler *Handler) {
	handler.Register(router)
}

func canonicalOrigin(raw string) (string, error) {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("invalid MCP OAuth site base URL")
	}
	return (&url.URL{Scheme: parsed.Scheme, Host: parsed.Host}).String(), nil
}
