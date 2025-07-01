package buscore

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type QueueResponse[T any] struct {
	Data T
}

type QueueSubscribeCallback[Req, Res any] func(ctx context.Context, data Req) Res

type Queue[Req, Res any] interface {
	Publish(ctx context.Context, data Req) error
	Request(ctx context.Context, data Req) (*QueueResponse[Res], error)
	SubscribeGroup(queueGroup string, data QueueSubscribeCallback[Req, Res]) error
	Subscribe(data QueueSubscribeCallback[Req, Res]) error
	Unsubscribe()
}

type NatsQueue[Req, Res any] struct {
	nc           *nats.Conn
	subscription *nats.Subscription
	subject      string
	timeout      time.Duration
	encoder      QueueEncoder
}

func NewNatsQueue[Req, Res any](
	nc *nats.Conn,
	subject string,
	timeout time.Duration,
	encoder QueueEncoder,
) *NatsQueue[Req, Res] {
	return &NatsQueue[Req, Res]{
		nc:      nc,
		subject: subject,
		timeout: timeout,
		encoder: encoder,
	}
}

func (c *NatsQueue[Req, Res]) Request(ctx context.Context, req Req) (*QueueResponse[Res], error) {
	tracer := otel.Tracer("nats-publisher")
	ctx, span := tracer.Start(ctx, "Publish "+c.subject, trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	span.SetAttributes(
		attribute.String("messaging.system", "nats"),
		attribute.String("messaging.destination", c.subject),
	)

	reqBytes, err := natsEncode(c.encoder, req)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	replyKey := nats.NewInbox()

	msg := &nats.Msg{
		Subject: c.subject,
		Reply:   replyKey,
		Header:  nats.Header{},
		Data:    reqBytes,
		Sub:     nil,
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(msg.Header))

	resp, err := c.nc.RequestMsgWithContext(ctx, msg)
	if err != nil {
		return nil, err
	}

	res, err := natsDecode[Res](c.encoder, resp.Data)
	if err != nil {
		return nil, err
	}

	return &QueueResponse[Res]{
		Data: res,
	}, nil
}

func (c *NatsQueue[Req, Res]) SubscribeGroup(
	queueGroup string,
	cb QueueSubscribeCallback[Req, Res],
) error {
	tracer := otel.Tracer("nats-subscriber")

	sub, err := c.nc.QueueSubscribe(
		c.subject,
		queueGroup,
		func(m *nats.Msg) {
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
				defer cancel()
				ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(m.Header))

				newCtx, span := tracer.Start(
					ctx,
					"Read "+c.subject,
					trace.WithSpanKind(trace.SpanKindConsumer),
				)
				defer span.End()

				span.SetAttributes(
					attribute.String("messaging.system", "nats"),
					attribute.String("messaging.destination", c.subject),
					attribute.String("messaging.group", queueGroup),
				)

				data, err := natsDecode[Req](c.encoder, m.Data)
				if err != nil {
					return
				}

				response := cb(newCtx, data)

				responseBytes, err := natsEncode(c.encoder, response)
				if err != nil {
					return
				}

				c.nc.Publish(m.Reply, responseBytes)
			}()
		},
	)

	c.subscription = sub

	return err
}

func (c *NatsQueue[Req, Res]) Subscribe(
	cb QueueSubscribeCallback[Req, Res],
) error {
	tracer := otel.Tracer("nats-subscriber")

	sub, err := c.nc.Subscribe(
		c.subject,
		func(m *nats.Msg) {
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
				defer cancel()
				ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(m.Header))

				ctx, span := tracer.Start(
					ctx,
					"Read "+c.subject,
					trace.WithSpanKind(trace.SpanKindConsumer),
				)
				defer span.End()

				span.SetAttributes(
					attribute.String("messaging.system", "nats"),
					attribute.String("messaging.destination", c.subject),
				)

				data, err := natsDecode[Req](c.encoder, m.Data)
				if err != nil {
					return
				}

				response := cb(ctx, data)

				responseBytes, err := natsEncode(c.encoder, response)
				if err != nil {
					return
				}

				c.nc.Publish(m.Reply, responseBytes)
			}()
		},
	)

	c.subscription = sub

	return err
}

func (c *NatsQueue[Req, Res]) Publish(ctx context.Context, data Req) error {
	tracer := otel.Tracer("nats-publisher")
	ctx, span := tracer.Start(
		ctx,
		"Publish "+c.subject,
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer span.End()

	span.SetAttributes(
		attribute.String("messaging.system", "nats"),
		attribute.String("messaging.destination", c.subject),
	)

	dataBytes, err := natsEncode(c.encoder, data)
	if err != nil {
		return err
	}

	msg := &nats.Msg{
		Subject: c.subject,
		Header:  nats.Header{},
		Data:    dataBytes,
		Sub:     nil,
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(msg.Header))

	return c.nc.PublishMsg(msg)
}

func (c *NatsQueue[Req, Res]) Unsubscribe() {
	if c.subscription != nil {
		c.subscription.Unsubscribe()
	}
}
