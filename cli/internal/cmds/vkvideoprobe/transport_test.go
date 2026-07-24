package vkvideoprobe

import "testing"

func TestParseTransportSpec_acceptsOnlyVerifiedShape(t *testing.T) {
	// Given
	spec := []byte(`{"url":"wss://pubsub.example.test/connection?protocol=v2","subprotocols":["json"],"frames":[{"connect":"{{connection_token}}"},{"channel":"{{channel}}","token":"{{subscription_token}}"}]}`)

	// When
	got, err := ParseTransportSpec(spec)

	// Then
	if err != nil {
		t.Fatalf("parse transport spec: %v", err)
	}
	if got.URL != "wss://pubsub.example.test/connection?protocol=v2" || len(got.Frames) != 2 {
		t.Fatalf("parsed spec = %#v", got)
	}
}

func TestParseTransportSpec_rejectsUnknownPlaceholders(t *testing.T) {
	// Given
	spec := []byte(`{"url":"wss://pubsub.example.test/connection","frames":[{"token":"{{unexpected}}"}]}`)

	// When
	_, err := ParseTransportSpec(spec)

	// Then
	if err == nil {
		t.Fatal("expected transport spec error")
	}
}

func TestParseTransportSpec_rejectsStaticSensitiveValues(t *testing.T) {
	// Given
	spec := []byte(`{"url":"wss://pubsub.example.test/connection","frames":[{"token":"synthetic-static-token"}]}`)

	// When
	_, err := ParseTransportSpec(spec)

	// Then
	if err == nil {
		t.Fatal("expected transport spec error")
	}
}
