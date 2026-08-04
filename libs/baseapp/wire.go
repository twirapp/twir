//go:build wireinject

package baseapp

import "github.com/goforj/wire"

func NewBase(opts Opts) (Base, error) {
	wire.Build(providerSet, wire.Struct(new(Base), "*"))
	return Base{}, nil
}
