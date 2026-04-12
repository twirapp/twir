package auth

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	sessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpdelivery "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	kickplatform "github.com/twirapp/twir/apps/api-gql/internal/platform/kick"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/repositories/tokens"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Huma             huma.API
	Gorm             *gorm.DB
	Config           config.Config
	Bus              *buscore.Bus
	Sessions         *sessions.Auth
	TokensRepository tokens.Repository
	KickProvider     *kickplatform.Provider
}

type Auth struct {
	gorm             *gorm.DB
	config           config.Config
	bus              *buscore.Bus
	sessions         *sessions.Auth
	tokensRepository tokens.Repository
	kickProvider     *kickplatform.Provider
}

func New(opts Opts) *Auth {
	p := &Auth{
		gorm:             opts.Gorm,
		config:           opts.Config,
		bus:              opts.Bus,
		sessions:         opts.Sessions,
		tokensRepository: opts.TokensRepository,
		kickProvider:     opts.KickProvider,
	}

	huma.Register(
		opts.Huma,
		huma.Operation{
			OperationID: "auth-post-code",
			Method:      http.MethodPost,
			Path:        "/auth",
			Tags:        []string{"Auth"},
			Summary:     "Auth post code",
		},
		func(
			ctx context.Context, i *struct {
				Body authBody
			},
		) (*httpdelivery.BaseOutputJson[authResponseDto], error) {
			return p.handleAuthPostCode(ctx, i.Body)
		},
	)

	huma.Register(
		opts.Huma,
		huma.Operation{
			OperationID: "auth-kick-authorize",
			Method:      http.MethodGet,
			Path:        "/auth/kick/authorize",
			Tags:        []string{"Auth"},
			Summary:     "Get Kick OAuth authorize URL",
		},
		func(ctx context.Context, _ *struct{}) (*kickAuthorizeOutput, error) {
			return p.handleKickAuthorize(ctx)
		},
	)

	return p
}
