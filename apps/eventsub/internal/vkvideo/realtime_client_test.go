package vkvideo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"

	centrifuge "github.com/centrifugal/centrifuge-go"
	"github.com/google/uuid"
)

type lifecycleLogEntry struct {
	Level         string `json:"level"`
	Message       string `json:"msg"`
	BindingID     string `json:"binding_id"`
	Code          uint64 `json:"code"`
	Reason        string `json:"reason"`
	Version       string `json:"version"`
	WasRecovering bool   `json:"was_recovering"`
	Recovered     bool   `json:"recovered"`
	Error         string `json:"error"`
}

func TestRealtimeClient_handlePublicationLogsReceiptMetadata(t *testing.T) {
	// Given
	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuffer, &slog.HandlerOptions{}))
	session, err := NewPublicationSession(1)
	if err != nil {
		t.Fatalf("create publication session: %v", err)
	}
	client := &RealtimeClient{
		session:   session,
		bindingID: uuid.New(),
		logger:    logger,
	}
	publication := []byte("fixture-publication")

	// When
	client.handlePublication(publication)
	queued, err := session.Receive(context.Background())

	// Then
	if err != nil {
		t.Fatalf("receive queued publication: %v", err)
	}
	if string(queued) != string(publication) {
		t.Fatalf("queued publication = %q, want %q", queued, publication)
	}
	if strings.Contains(logBuffer.String(), string(publication)) {
		t.Fatal("publication payload leaked into log output")
	}

	var entry map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(logBuffer.Bytes()), &entry); err != nil {
		t.Fatalf("decode log entry: %v", err)
	}
	if got := entry["level"]; got != "INFO" {
		t.Fatalf("log level = %v, want INFO", got)
	}
	if got := entry["binding_id"]; got != client.bindingID.String() {
		t.Fatalf("binding_id = %v, want %s", got, client.bindingID)
	}
	if got := entry["publication_size"]; got != float64(len(publication)) {
		t.Fatalf("publication_size = %v, want %d", got, len(publication))
	}
	if got := entry["enqueue_result"]; got != true {
		t.Fatalf("enqueue_result = %v, want true", got)
	}
}

func TestRealtimeClient_logsLifecycleStateSafely(t *testing.T) {
	// Given
	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuffer, &slog.HandlerOptions{}))
	client := &RealtimeClient{
		bindingID: uuid.New(),
		logger:    logger,
	}

	// When
	client.logConnecting(centrifuge.ConnectingEvent{Code: 11, Reason: "connect called"})
	client.logConnected(centrifuge.ConnectedEvent{Version: "v2"})
	client.logDisconnected(centrifuge.DisconnectedEvent{Code: 12, Reason: "transport closed"})
	client.logConnectionError(centrifuge.ErrorEvent{Error: errors.New("connection failed")})
	client.logSubscribing(centrifuge.SubscribingEvent{Code: 13, Reason: "subscribe called"})
	client.logSubscribed(centrifuge.SubscribedEvent{Recovered: true, WasRecovering: false})
	client.logSubscriptionError(centrifuge.SubscriptionErrorEvent{Error: errors.New("subscription failed")})

	// Then
	lines := strings.Split(strings.TrimSpace(logBuffer.String()), "\n")
	if len(lines) != 7 {
		t.Fatalf("log line count = %d, want 7", len(lines))
	}

	entries := make([]lifecycleLogEntry, 0, len(lines))
	for _, line := range lines {
		var entry lifecycleLogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatalf("decode log entry: %v", err)
		}
		entries = append(entries, entry)
	}

	if got := entries[0]; got.Level != "INFO" || got.Message != "VK Video Centrifugo connecting" || got.BindingID != client.bindingID.String() || got.Code != 11 || got.Reason != "connect called" {
		t.Fatalf("connecting log = %+v", got)
	}
	if got := entries[1]; got.Level != "INFO" || got.Message != "VK Video Centrifugo connected" || got.BindingID != client.bindingID.String() || got.Version != "v2" {
		t.Fatalf("connected log = %+v", got)
	}
	if got := entries[2]; got.Level != "WARN" || got.Message != "VK Video Centrifugo disconnected" || got.BindingID != client.bindingID.String() || got.Code != 12 || got.Reason != "transport closed" {
		t.Fatalf("disconnected log = %+v", got)
	}
	if got := entries[3]; got.Level != "ERROR" || got.Message != "VK Video Centrifugo connection error" || got.BindingID != client.bindingID.String() || got.Error != "connection failed" {
		t.Fatalf("connection error log = %+v", got)
	}
	if got := entries[4]; got.Level != "INFO" || got.Message != "VK Video Centrifugo subscribing" || got.BindingID != client.bindingID.String() || got.Code != 13 || got.Reason != "subscribe called" {
		t.Fatalf("subscribing log = %+v", got)
	}
	if got := entries[5]; got.Level != "INFO" || got.Message != "VK Video Centrifugo subscribed" || got.BindingID != client.bindingID.String() || !got.Recovered || got.WasRecovering {
		t.Fatalf("subscribed log = %+v", got)
	}
	if got := entries[6]; got.Level != "ERROR" || got.Message != "VK Video Centrifugo subscription error" || got.BindingID != client.bindingID.String() || got.Error != "subscription failed" {
		t.Fatalf("subscription error log = %+v", got)
	}

	if strings.Contains(logBuffer.String(), centrifugoEndpoint) {
		t.Fatal("endpoint leaked into lifecycle log output")
	}
	if strings.Contains(logBuffer.String(), "channel-chat:") {
		t.Fatal("channel leaked into lifecycle log output")
	}
	if strings.Contains(logBuffer.String(), "fixture-publication") {
		t.Fatal("publication payload leaked into lifecycle log output")
	}
}
