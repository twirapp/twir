package wsrouter

import (
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
)

func NewNatsSubscription(opts Opts) (*WsRouterNats, error) {
	nc, err := nats.Connect(opts.Config.NatsUrl)
	if err != nil {
		return nil, err
	}

	return &WsRouterNats{
		nc: nc,
	}, nil
}

type WsRouterNats struct {
	nc *nats.Conn
}

var _ WsRouter = &WsRouterNats{}

type WsRouterNatsSubscription struct {
	subs      []*nats.Subscription
	dataChann chan []byte
}

func (c *WsRouterNatsSubscription) Unsubscribe() error {
	for _, sub := range c.subs {
		if err := sub.Unsubscribe(); err != nil {
			return err
		}
	}

	return nil
}

func (c *WsRouterNatsSubscription) GetChannel() chan []byte {
	return c.dataChann
}

func (c *WsRouterNats) Subscribe(keys []string) (WsRouterSubscription, error) {
	ch := make(chan []byte)
	subs := make([]*nats.Subscription, 0, len(keys))

	for _, key := range keys {
		sub, err := c.nc.Subscribe(
			key,
			func(msg *nats.Msg) {
				ch <- msg.Data
			},
		)

		if err != nil {
			return nil, err
		}

		subs = append(subs, sub)
	}

	return &WsRouterNatsSubscription{
		subs:      subs,
		dataChann: ch,
	}, nil
}

func (c *WsRouterNats) Publish(key string, data any) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.nc.Publish(key, dataBytes)
}
