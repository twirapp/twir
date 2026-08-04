//go:build wireinject

package baseapp

import "github.com/goforj/wire"

func NewBase(opts Opts) (Base, error) {
	wire.Build(providerSet, newBase)
	return Base{}, nil
}
