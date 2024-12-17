package commands_responses

import (
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
}

func New(opts Opts) *Service {
	return &Service{}
}

type Service struct {
}
