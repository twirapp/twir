package wsrouter

import (
	config "github.com/twirapp/twir/libs/config"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Config config.Config
}

type WsRouter interface {
	Subscribe(keys []string) (WsRouterSubscription, error)
	Publish(key string, data any) error
}

type WsRouterSubscription interface {
	GetChannel() chan []byte
	Unsubscribe() error
}
