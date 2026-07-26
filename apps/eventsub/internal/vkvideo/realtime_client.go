package vkvideo

import (
	"errors"
	"fmt"
	"log/slog"

	centrifuge "github.com/centrifugal/centrifuge-go"
	"github.com/google/uuid"
)

var ErrInvalidRealtimeClientConfig = errors.New("vk video realtime client configuration is invalid")

const centrifugoEndpoint = "wss://pubsub-dev.live.vkvideo.ru/connection/websocket?format=json&cf_protocol_version=v2"

type TokenCallbacks struct {
	Connection   func() (string, error)
	Subscription func(channel string) (string, error)
}

type RealtimeClientConfig struct {
	Channel       string
	QueueCapacity int
	BindingID     uuid.UUID
	Logger        *slog.Logger
	Tokens        TokenCallbacks
}

type RealtimeClient struct {
	client    *centrifuge.Client
	session   *PublicationSession
	sub       *centrifuge.Subscription
	bindingID uuid.UUID
	logger    *slog.Logger
}

func NewRealtimeClient(config RealtimeClientConfig) (*RealtimeClient, error) {
	if config.Channel == "" {
		return nil, fmt.Errorf("channel is required: %w", ErrInvalidRealtimeClientConfig)
	}
	if config.Tokens.Connection == nil || config.Tokens.Subscription == nil {
		return nil, fmt.Errorf("connection and subscription token callbacks are required: %w", ErrInvalidRealtimeClientConfig)
	}

	session, err := NewPublicationSession(config.QueueCapacity)
	if err != nil {
		return nil, fmt.Errorf("create publication session: %w", err)
	}

	realtimeClient := &RealtimeClient{
		session:   session,
		bindingID: config.BindingID,
		logger:    config.Logger,
	}
	centrifugoClient := centrifuge.NewJsonClient(centrifugoEndpoint, centrifuge.Config{
		GetToken: func(centrifuge.ConnectionTokenEvent) (string, error) {
			return config.Tokens.Connection()
		},
	})
	subscription, err := centrifugoClient.NewSubscription(config.Channel, centrifuge.SubscriptionConfig{
		GetToken: func(event centrifuge.SubscriptionTokenEvent) (string, error) {
			return config.Tokens.Subscription(event.Channel)
		},
	})
	if err != nil {
		centrifugoClient.Close()
		return nil, fmt.Errorf("create Centrifugo subscription: %w", err)
	}

	realtimeClient.client = centrifugoClient
	realtimeClient.sub = subscription
	realtimeClient.registerLifecycleCallbacks()

	subscription.OnPublication(func(event centrifuge.PublicationEvent) {
		realtimeClient.handlePublication(event.Data)
	})

	return realtimeClient, nil
}

func (c *RealtimeClient) registerLifecycleCallbacks() {
	if c.client != nil {
		c.client.OnConnecting(func(event centrifuge.ConnectingEvent) {
			c.logConnecting(event)
		})
		c.client.OnConnected(func(event centrifuge.ConnectedEvent) {
			c.logConnected(event)
		})
		c.client.OnDisconnected(func(event centrifuge.DisconnectedEvent) {
			c.logDisconnected(event)
		})
		c.client.OnError(func(event centrifuge.ErrorEvent) {
			c.logConnectionError(event)
		})
	}

	if c.sub != nil {
		c.sub.OnSubscribing(func(event centrifuge.SubscribingEvent) {
			c.logSubscribing(event)
		})
		c.sub.OnSubscribed(func(event centrifuge.SubscribedEvent) {
			c.logSubscribed(event)
		})
		c.sub.OnError(func(event centrifuge.SubscriptionErrorEvent) {
			c.logSubscriptionError(event)
		})
	}
}

func (c *RealtimeClient) logConnecting(event centrifuge.ConnectingEvent) {
	if c.logger == nil {
		return
	}

	c.logger.Info("VK Video Centrifugo connecting",
		slog.String("binding_id", c.bindingID.String()),
		slog.Uint64("code", uint64(event.Code)),
		slog.String("reason", event.Reason),
	)
}

func (c *RealtimeClient) logConnected(event centrifuge.ConnectedEvent) {
	if c.logger == nil {
		return
	}

	c.logger.Info("VK Video Centrifugo connected",
		slog.String("binding_id", c.bindingID.String()),
		slog.String("version", event.Version),
	)
}

func (c *RealtimeClient) logDisconnected(event centrifuge.DisconnectedEvent) {
	if c.logger == nil {
		return
	}

	c.logger.Warn("VK Video Centrifugo disconnected",
		slog.String("binding_id", c.bindingID.String()),
		slog.Uint64("code", uint64(event.Code)),
		slog.String("reason", event.Reason),
	)
}

func (c *RealtimeClient) logConnectionError(event centrifuge.ErrorEvent) {
	if c.logger == nil {
		return
	}

	c.logger.Error("VK Video Centrifugo connection error",
		slog.String("binding_id", c.bindingID.String()),
		slog.String("error", event.Error.Error()),
	)
}

func (c *RealtimeClient) logSubscribing(event centrifuge.SubscribingEvent) {
	if c.logger == nil {
		return
	}

	c.logger.Info("VK Video Centrifugo subscribing",
		slog.String("binding_id", c.bindingID.String()),
		slog.Uint64("code", uint64(event.Code)),
		slog.String("reason", event.Reason),
	)
}

func (c *RealtimeClient) logSubscribed(event centrifuge.SubscribedEvent) {
	if c.logger == nil {
		return
	}

	c.logger.Info("VK Video Centrifugo subscribed",
		slog.String("binding_id", c.bindingID.String()),
		slog.Bool("was_recovering", event.WasRecovering),
		slog.Bool("recovered", event.Recovered),
	)
}

func (c *RealtimeClient) logSubscriptionError(event centrifuge.SubscriptionErrorEvent) {
	if c.logger == nil {
		return
	}

	c.logger.Error("VK Video Centrifugo subscription error",
		slog.String("binding_id", c.bindingID.String()),
		slog.String("error", event.Error.Error()),
	)
}

func (c *RealtimeClient) handlePublication(publication []byte) {
	enqueued := c.session.Enqueue(publication)
	if c.logger == nil {
		return
	}

	c.logger.Info("VK Video Centrifugo publication received",
		slog.String("binding_id", c.bindingID.String()),
		slog.Int("publication_size", len(publication)),
		slog.Bool("enqueue_result", enqueued),
	)
}

func (c *RealtimeClient) Connect() error {
	if err := c.sub.Subscribe(); err != nil {
		return fmt.Errorf("subscribe to Centrifugo channel: %w", err)
	}
	if err := c.client.Connect(); err != nil {
		return fmt.Errorf("connect to Centrifugo: %w", err)
	}

	return nil
}

func (c *RealtimeClient) Session() *PublicationSession {
	return c.session
}

func (c *RealtimeClient) Close() {
	c.session.Close()
	c.client.Close()
}
