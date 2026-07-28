package vkvideo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	centrifuge "github.com/centrifugal/centrifuge-go"
	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/integrations/vk"
)

var (
	ErrInvalidRealtimeClientConfig = errors.New("vk video realtime client configuration is invalid")
	ErrEmptySubscriptionToken      = errors.New("vk video subscription token is empty")
)

const centrifugoEndpoint = "wss://pubsub-dev.live.vkvideo.ru/connection/websocket?format=json&cf_protocol_version=v2"

type TokenCallbacks struct {
	Context      context.Context
	Connection   func(context.Context) (string, error)
	Subscription func(context.Context, string) (string, error)
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
	channel   string
	tokens    TokenCallbacks
	bindingID uuid.UUID
	logger    *slog.Logger
}

func NewRealtimeClient(config RealtimeClientConfig) (*RealtimeClient, error) {
	return newRealtimeClient(config, centrifugoEndpoint)
}

func newRealtimeClient(config RealtimeClientConfig, endpoint string) (*RealtimeClient, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint is required: %w", ErrInvalidRealtimeClientConfig)
	}
	if config.Channel == "" {
		return nil, fmt.Errorf("channel is required: %w", ErrInvalidRealtimeClientConfig)
	}
	if config.Tokens.Context == nil {
		return nil, fmt.Errorf("token context is required: %w", ErrInvalidRealtimeClientConfig)
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
		channel:   config.Channel,
		tokens:    config.Tokens,
		bindingID: config.BindingID,
		logger:    config.Logger,
	}
	realtimeClient.client = centrifuge.NewJsonClient(endpoint, centrifuge.Config{
		GetToken: func(centrifuge.ConnectionTokenEvent) (string, error) {
			return realtimeClient.tokens.Connection(realtimeClient.tokens.Context)
		},
	})

	return realtimeClient, nil
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

func (c *RealtimeClient) Connect(ctx context.Context) error {
	signals := newRealtimeStartupSignals()
	c.registerClientLifecycleCallbacks(signals)

	if err := c.client.Connect(); err != nil {
		return fmt.Errorf("connect to Centrifugo: %w", err)
	}
	if err := signals.await(ctx, signals.connected); err != nil {
		return fmt.Errorf("wait for Centrifugo connection: %w", err)
	}

	token, err := c.tokens.Subscription(ctx, c.channel)
	subscriptionConfig := centrifuge.SubscriptionConfig{}
	if err != nil && !errors.Is(err, vk.ErrWebSocketSubscriptionChannelTokenMissing) {
		return fmt.Errorf("get Centrifugo subscription token: %w", err)
	}
	if err == nil && token == "" {
		return ErrEmptySubscriptionToken
	}
	if err == nil {
		subscriptionConfig = centrifuge.SubscriptionConfig{
			Token: token,
			GetToken: func(event centrifuge.SubscriptionTokenEvent) (string, error) {
				return c.tokens.Subscription(c.tokens.Context, event.Channel)
			},
		}
	}

	subscription, err := c.client.NewSubscription(c.channel, subscriptionConfig)
	if err != nil {
		return fmt.Errorf("create Centrifugo subscription: %w", err)
	}
	c.registerSubscriptionLifecycleCallbacks(subscription, signals)

	if err := subscription.Subscribe(); err != nil {
		return fmt.Errorf("subscribe to Centrifugo channel: %w", err)
	}
	if err := signals.await(ctx, signals.subscribed); err != nil {
		return fmt.Errorf("wait for Centrifugo subscription: %w", err)
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
