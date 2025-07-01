package auth

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	config "github.com/satont/twir/libs/config"
	sessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpdelivery "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Huma     huma.API
	Gorm     *gorm.DB
	Config   config.Config
	Bus      *buscore.Bus
	Sessions *sessions.Auth
}

type Auth struct {
	gorm     *gorm.DB
	config   config.Config
	bus      *buscore.Bus
	sessions *sessions.Auth
}

func New(opts Opts) *Auth {
	p := &Auth{
		gorm:     opts.Gorm,
		config:   opts.Config,
		bus:      opts.Bus,
		sessions: opts.Sessions,
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

	return p
}
