package vkvideo

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	centrifuge "github.com/centrifugal/centrifuge-go"
)

type realtimeStartupSignals struct {
	connected  chan struct{}
	subscribed chan struct{}
	failure    chan error

	connectedOnce  sync.Once
	subscribedOnce sync.Once
	failureOnce    sync.Once
}

func newRealtimeStartupSignals() *realtimeStartupSignals {
	return &realtimeStartupSignals{
		connected:  make(chan struct{}, 1),
		subscribed: make(chan struct{}, 1),
		failure:    make(chan error, 1),
	}
}

func (s *realtimeStartupSignals) signalConnected() {
	s.connectedOnce.Do(func() {
		select {
		case s.connected <- struct{}{}:
		default:
		}
	})
}

func (s *realtimeStartupSignals) signalSubscribed() {
	s.subscribedOnce.Do(func() {
		select {
		case s.subscribed <- struct{}{}:
		default:
		}
	})
}

func (s *realtimeStartupSignals) signalFailure(err error) {
	s.failureOnce.Do(func() {
		select {
		case s.failure <- err:
		default:
		}
	})
}

func (s *realtimeStartupSignals) await(ctx context.Context, ready <-chan struct{}) error {
	select {
	case <-ready:
		return nil
	case err := <-s.failure:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *RealtimeClient) registerClientLifecycleCallbacks(signals *realtimeStartupSignals) {
	c.client.OnConnecting(func(event centrifuge.ConnectingEvent) {
		c.logConnecting(event)
	})
	c.client.OnConnected(func(event centrifuge.ConnectedEvent) {
		c.logConnected(event)
		signals.signalConnected()
	})
	c.client.OnDisconnected(func(event centrifuge.DisconnectedEvent) {
		c.logDisconnected(event)
		signals.signalFailure(fmt.Errorf(
			"Centrifugo disconnected during startup with code %d: %s",
			event.Code,
			event.Reason,
		))
	})
	c.client.OnError(func(event centrifuge.ErrorEvent) {
		c.logConnectionError(event)
		signals.signalFailure(event.Error)
	})
}

func (c *RealtimeClient) registerSubscriptionLifecycleCallbacks(
	subscription *centrifuge.Subscription,
	signals *realtimeStartupSignals,
) {
	subscription.OnSubscribing(func(event centrifuge.SubscribingEvent) {
		c.logSubscribing(event)
	})
	subscription.OnSubscribed(func(event centrifuge.SubscribedEvent) {
		c.logSubscribed(event)
		signals.signalSubscribed()
	})
	subscription.OnError(func(event centrifuge.SubscriptionErrorEvent) {
		c.logSubscriptionError(event)
		signals.signalFailure(event.Error)
	})
	subscription.OnUnsubscribed(func(event centrifuge.UnsubscribedEvent) {
		signals.signalFailure(fmt.Errorf(
			"Centrifugo subscription ended during startup with code %d: %s",
			event.Code,
			event.Reason,
		))
	})
	subscription.OnPublication(func(event centrifuge.PublicationEvent) {
		c.handlePublication(event.Data)
	})
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
