package resolvers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/sessions"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	sessions *sessions.Sessions
	gorm     *gorm.DB

	clientsCommandsChannels map[string]chan *gqlmodel.Command
}

type Opts struct {
	fx.In

	Sessions *sessions.Sessions
	Gorm     *gorm.DB
}

func New(opts Opts) *Resolver {
	return &Resolver{
		sessions:                opts.Sessions,
		gorm:                    opts.Gorm,
		clientsCommandsChannels: make(map[string]chan *gqlmodel.Command, 1),
	}
}
