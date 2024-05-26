package wsrouter

import (
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
	config "github.com/satont/twir/libs/config"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Config config.Config
}

func New(opts Opts) (*WsSync, error) {
	nc, err := nats.Connect(opts.Config.NatsUrl)
	if err != nil {
		return nil, err
	}

	return &WsSync{
		nc: nc,
	}, nil
}

type WsSync struct {
	nc *nats.Conn
}

type Subscription struct {
	subs []*nats.Subscription
	Data chan []byte
}

func (c *Subscription) Unsubscribe() error {
	for _, sub := range c.subs {
		if err := sub.Unsubscribe(); err != nil {
			return err
		}
	}

	return nil
}

func (c *WsSync) Subscribe(keys []string) (*Subscription, error) {
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

	return &Subscription{
		subs: subs,
		Data: ch,
	}, nil
}

func (c *WsSync) Publish(key string, data any) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.nc.Publish(key, dataBytes)
}
