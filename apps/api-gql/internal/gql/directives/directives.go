package directives

import (
	"github.com/twirapp/twir/apps/api-gql/internal/sessions"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Sessions *sessions.Sessions
}

func New(opts Opts) *Directives {
	return &Directives{
		sessions: opts.Sessions,
	}
}

type Directives struct {
	sessions *sessions.Sessions
}
