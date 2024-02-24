package services

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
)

type QueueResponse[T any] struct {
	Data T
}

type QueueSubscribeCallback[Req, Res any] func(ctx context.Context, data Req) Res

type Queue[Req, Res any] interface {
	Publish(data Req) error
	Request(ctx context.Context, data Req) (*QueueResponse[Res], error)
	Subscribe(data QueueSubscribeCallback[Req, Res]) error
	Unsubscribe()
}

type NatsQueue[Req, Res any] struct {
	nc           *nats.EncodedConn
	queue        string
	timeout      time.Duration
	subscription *nats.Subscription
}

func NewNatsQueue[Req, Res any](
	nc *nats.Conn,
	queue string,
	timeout time.Duration,
) *NatsQueue[Req, Res] {
	newNc, _ := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	return &NatsQueue[Req, Res]{
		nc:      newNc,
		queue:   queue,
		timeout: timeout,
	}
}

func (c *NatsQueue[Req, Res]) Request(ctx context.Context, req Req) (*QueueResponse[Res], error) {
	var res Res

	err := c.nc.RequestWithContext(ctx, c.queue, &req, &res)
	if err != nil {
		return nil, err
	}

	return &QueueResponse[Res]{
		Data: res,
	}, nil
}

func (c *NatsQueue[Req, Res]) Subscribe(
	cb QueueSubscribeCallback[Req, Res],
) error {
	sub, err := c.nc.Subscribe(
		c.queue,
		func(subject, reply string, data *Req) {
			ctx, _ := context.WithTimeout(context.Background(), c.timeout)
			response := cb(ctx, *data)
			c.nc.Publish(reply, &response)
		},
	)

	c.subscription = sub

	return err
}

func (c *NatsQueue[Req, Res]) Publish(data Req) error {
	return c.nc.Publish(c.queue, &data)
}

func (c *NatsQueue[Req, Res]) Unsubscribe() {
	if c.subscription != nil {
		c.subscription.Unsubscribe()
	}
}
