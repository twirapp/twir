package wsrouter

import (
	"fmt"
	"sync"

	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
	config "github.com/twirapp/twir/libs/config"
	"go.uber.org/fx"
)

const subscriptionBufferSize = 64

func NewNatsSubscription(opts Opts) (*WsRouterNats, error) {
	nc, err := nats.Connect(opts.Config.NatsUrl)
	if err != nil {
		return nil, err
	}

	return &WsRouterNats{
		nc: nc,
	}, nil
}

func NewNatsWsRouterFx(cfg config.Config) (*WsRouterNats, error) {
	nc, err := nats.Connect(cfg.NatsUrl)
	if err != nil {
		return nil, err
	}

	return &WsRouterNats{
		nc: nc,
	}, nil
}

type NatsWsRouterFxOpts struct {
	fx.In

	Config config.Config
}

func NewNatsWsRouterFxWithOpts(opts NatsWsRouterFxOpts) (*WsRouterNats, error) {
	return NewNatsWsRouterFx(opts.Config)
}

type WsRouterNats struct {
	nc *nats.Conn
}

var _ WsRouter = &WsRouterNats{}

type WsRouterNatsSubscription struct {
	subs           []*nats.Subscription
	dataChann      chan []byte
	done           chan struct{}
	unsubscribeErr error
	unsubscribe    sync.Once
}

func (c *WsRouterNatsSubscription) Unsubscribe() error {
	c.unsubscribe.Do(func() {
		close(c.done)
		for _, sub := range c.subs {
			if err := sub.Unsubscribe(); err != nil && c.unsubscribeErr == nil {
				c.unsubscribeErr = err
			}
		}
	})

	return c.unsubscribeErr
}

func (c *WsRouterNatsSubscription) GetChannel() chan []byte {
	return c.dataChann
}

func (c *WsRouterNats) Subscribe(keys []string) (WsRouterSubscription, error) {
	subscription := &WsRouterNatsSubscription{
		subs:      make([]*nats.Subscription, 0, len(keys)),
		dataChann: make(chan []byte, subscriptionBufferSize),
		done:      make(chan struct{}),
	}

	for _, key := range keys {
		sub, err := c.nc.Subscribe(
			key,
			func(msg *nats.Msg) {
				select {
				case subscription.dataChann <- msg.Data:
				case <-subscription.done:
				}
			},
		)
		if err != nil {
			_ = subscription.Unsubscribe()
			return nil, fmt.Errorf("subscribe to NATS key %q: %w", key, err)
		}

		subscription.subs = append(subscription.subs, sub)
	}

	if err := c.nc.Flush(); err != nil {
		_ = subscription.Unsubscribe()
		return nil, fmt.Errorf("flush NATS subscriptions: %w", err)
	}

	return subscription, nil
}

func (c *WsRouterNats) Publish(key string, data any) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.nc.Publish(key, dataBytes)
}
