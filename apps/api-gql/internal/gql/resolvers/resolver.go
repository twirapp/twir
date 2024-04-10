package resolvers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/sessions"
	"go.uber.org/fx"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	sessions        *sessions.Sessions
	NewCommandChann chan *gqlmodel.Command
}

type Opts struct {
	fx.In

	Sessions *sessions.Sessions
}

func New(opts Opts) *Resolver {
	return &Resolver{
		sessions:        opts.Sessions,
		NewCommandChann: make(chan *gqlmodel.Command),
	}
}
