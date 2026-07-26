package vkvideo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/twirapp/twir/libs/integrations/vk"
)

const (
	realtimeProtocolTestChannel           = "chat:test"
	realtimeProtocolTestConnectionToken   = "connection-token"
	realtimeProtocolTestSubscriptionToken = "subscription-token"
)

func TestRealtimeClientConnectWaitsForSubscriptionAcknowledgement(t *testing.T) {
	// Given
	server := newCentrifugoProtocolServer(t, realtimeProtocolTestChannel)
	client := newProtocolRealtimeClient(t, server.endpoint, func(context.Context, string) (string, error) {
		return realtimeProtocolTestSubscriptionToken, nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), protocolCommandTimeout)
	defer cancel()
	connectResult := make(chan error, 1)

	// When
	go func() {
		connectResult <- client.Connect(ctx)
	}()
	connect, err := server.AwaitConnect(ctx)
	if err != nil {
		t.Fatalf("wait for connect command: %v", err)
	}
	if connect.Token != realtimeProtocolTestConnectionToken {
		t.Fatal("connect command did not carry the connection token")
	}
	server.ReleaseConnectAcknowledgement()
	subscribe, err := server.AwaitSubscribe(ctx)
	if err != nil {
		t.Fatalf("wait for subscribe command: %v", err)
	}
	if subscribe.Channel != realtimeProtocolTestChannel {
		t.Fatal("subscribe command channel did not match the configured channel")
	}
	if subscribe.Token != realtimeProtocolTestSubscriptionToken {
		t.Fatal("subscribe command did not carry the subscription token")
	}

	// Then
	timer := time.NewTimer(150 * time.Millisecond)
	defer timer.Stop()
	select {
	case err := <-connectResult:
		t.Fatalf("Connect returned before subscribe acknowledgement: %v", err)
	case <-timer.C:
	}
	server.ReleaseSubscribeAcknowledgement()
	if err := awaitConnectResult(ctx, connectResult); err != nil {
		t.Fatalf("Connect after subscribe acknowledgement: %v", err)
	}
}

func TestRealtimeClientConnectReturnsSubscriptionTokenError(t *testing.T) {
	// Given
	sentinel := errors.New("subscription token callback failed")
	server := newCentrifugoProtocolServer(t, realtimeProtocolTestChannel)
	client := newProtocolRealtimeClient(t, server.endpoint, func(context.Context, string) (string, error) {
		return "", sentinel
	})
	ctx, cancel := context.WithTimeout(context.Background(), protocolCommandTimeout)
	defer cancel()
	connectResult := make(chan error, 1)

	// When
	go func() {
		connectResult <- client.Connect(ctx)
	}()
	if _, err := server.AwaitConnect(ctx); err != nil {
		t.Fatalf("wait for connect command: %v", err)
	}
	server.ReleaseConnectAcknowledgement()

	// Then
	err := awaitConnectResult(ctx, connectResult)
	if !errors.Is(err, sentinel) {
		t.Fatalf("Connect error = %v, want subscription token callback error", err)
	}
	server.RequireNoSubscribe(t)
}

func TestRealtimeClientConnectRejectsEmptySubscriptionToken(t *testing.T) {
	// Given
	server := newCentrifugoProtocolServer(t, realtimeProtocolTestChannel)
	client := newProtocolRealtimeClient(t, server.endpoint, func(context.Context, string) (string, error) {
		return "", nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), protocolCommandTimeout)
	defer cancel()
	connectResult := make(chan error, 1)

	// When
	go func() {
		connectResult <- client.Connect(ctx)
	}()
	if _, err := server.AwaitConnect(ctx); err != nil {
		t.Fatalf("wait for connect command: %v", err)
	}
	server.ReleaseConnectAcknowledgement()

	// Then
	err := awaitConnectResult(ctx, connectResult)
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "empty") || !strings.Contains(strings.ToLower(err.Error()), "token") {
		t.Fatalf("Connect error = %v, want error mentioning empty token", err)
	}
	server.RequireNoSubscribe(t)
}

func TestRealtimeClientConnectSubscribesWithoutTokenWhenChannelTokenMissing(t *testing.T) {
	// Given
	server := newCentrifugoProtocolServer(t, realtimeProtocolTestChannel)
	client := newProtocolRealtimeClient(t, server.endpoint, func(context.Context, string) (string, error) {
		return "", fmt.Errorf("subscription token callback: %w", vk.ErrWebSocketSubscriptionChannelTokenMissing)
	})
	ctx, cancel := context.WithTimeout(context.Background(), protocolCommandTimeout)
	defer cancel()
	connectResult := make(chan error, 1)
	payload := []byte(`{"event":"publication","sequence":1}`)

	// When
	go func() {
		connectResult <- client.Connect(ctx)
	}()
	if _, err := server.AwaitConnect(ctx); err != nil {
		t.Fatalf("wait for connect command: %v", err)
	}
	server.ReleaseConnectAcknowledgement()
	var subscribeToken string
	select {
	case command := <-server.subscribeCommands:
		subscribeToken = command.Subscribe.Token
	case err := <-connectResult:
		t.Fatalf("Connect error = %v, want tokenless public-channel subscription", err)
	case err := <-server.errors:
		t.Fatalf("protocol server error: %v", err)
	case <-ctx.Done():
		t.Fatalf("wait for subscribe command: %v", ctx.Err())
	}
	if subscribeToken != "" {
		t.Fatalf("subscribe command token = %q, want empty", subscribeToken)
	}
	server.ReleaseSubscribeAcknowledgement()
	if err := awaitConnectResult(ctx, connectResult); err != nil {
		t.Fatalf("Connect after subscribe acknowledgement: %v", err)
	}
	if err := server.SendPublication(payload); err != nil {
		t.Fatalf("send publication: %v", err)
	}

	// Then
	received, err := client.Session().Receive(ctx)
	if err != nil {
		t.Fatalf("receive publication: %v", err)
	}
	if !bytes.Equal(received, payload) {
		t.Fatal("received publication payload did not match the protocol push")
	}
}

func TestRealtimeClientConnectDeliversPublicationAfterSubscribed(t *testing.T) {
	// Given
	server := newCentrifugoProtocolServer(t, realtimeProtocolTestChannel)
	client := newProtocolRealtimeClient(t, server.endpoint, func(context.Context, string) (string, error) {
		return realtimeProtocolTestSubscriptionToken, nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), protocolCommandTimeout)
	defer cancel()
	connectResult := make(chan error, 1)
	payload := []byte(`{"event":"publication","sequence":1}`)

	// When
	go func() {
		connectResult <- client.Connect(ctx)
	}()
	connect, err := server.AwaitConnect(ctx)
	if err != nil {
		t.Fatalf("wait for connect command: %v", err)
	}
	if connect.Token != realtimeProtocolTestConnectionToken {
		t.Fatal("connect command did not carry the connection token")
	}
	server.ReleaseConnectAcknowledgement()
	subscribe, err := server.AwaitSubscribe(ctx)
	if err != nil {
		t.Fatalf("wait for subscribe command: %v", err)
	}
	if subscribe.Token != realtimeProtocolTestSubscriptionToken {
		t.Fatal("subscribe command did not carry the subscription token")
	}
	server.ReleaseSubscribeAcknowledgement()
	if err := awaitConnectResult(ctx, connectResult); err != nil {
		t.Fatalf("Connect after subscribe acknowledgement: %v", err)
	}
	if err := server.SendPublication(payload); err != nil {
		t.Fatalf("send publication: %v", err)
	}

	// Then
	received, err := client.Session().Receive(ctx)
	if err != nil {
		t.Fatalf("receive publication: %v", err)
	}
	if !bytes.Equal(received, payload) {
		t.Fatal("received publication payload did not match the protocol push")
	}
}

func newProtocolRealtimeClient(
	t *testing.T,
	endpoint string,
	subscriptionToken func(context.Context, string) (string, error),
) *RealtimeClient {
	t.Helper()

	client, err := newRealtimeClient(RealtimeClientConfig{
		Channel:       realtimeProtocolTestChannel,
		QueueCapacity: 1,
		Tokens: TokenCallbacks{
			Context: context.Background(),
			Connection: func(context.Context) (string, error) {
				return realtimeProtocolTestConnectionToken, nil
			},
			Subscription: subscriptionToken,
		},
	}, endpoint)
	if err != nil {
		t.Fatalf("create realtime client: %v", err)
	}
	t.Cleanup(client.Close)

	return client
}

func awaitConnectResult(ctx context.Context, results <-chan error) error {
	select {
	case err := <-results:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
