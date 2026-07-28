package vkvideo

import (
	"context"
	"errors"
	"testing"
)

func TestPublicationSession_EnqueueCopiesOpaqueBytes(t *testing.T) {
	// Given
	session, err := NewPublicationSession(1)
	if err != nil {
		t.Fatalf("create publication session: %v", err)
	}
	publication := []byte{0xff, 0x00, 0x01}

	// When
	accepted := session.Enqueue(publication)
	publication[0] = 0x02
	queued, err := session.Receive(context.Background())

	// Then
	if !accepted {
		t.Fatal("publication was not accepted")
	}
	if err != nil {
		t.Fatalf("receive queued publication: %v", err)
	}
	if string(queued) != string([]byte{0xff, 0x00, 0x01}) {
		t.Fatalf("queued publication = %v, want independent opaque copy", queued)
	}
}

func TestPublicationSession_EnqueueDropsWhenQueueIsSaturated(t *testing.T) {
	// Given
	session, err := NewPublicationSession(1)
	if err != nil {
		t.Fatalf("create publication session: %v", err)
	}
	if !session.Enqueue([]byte{1}) {
		t.Fatal("first publication was not accepted")
	}

	// When
	accepted := session.Enqueue([]byte{2})

	// Then
	if accepted {
		t.Fatal("saturated queue accepted publication")
	}
}

func TestPublicationSession_ReceiveStopsAfterClose(t *testing.T) {
	// Given
	session, err := NewPublicationSession(1)
	if err != nil {
		t.Fatalf("create publication session: %v", err)
	}
	session.Close()

	// When
	_, err = session.Receive(context.Background())

	// Then
	if !errors.Is(err, ErrPublicationSessionClosed) {
		t.Fatalf("receive error = %v, want ErrPublicationSessionClosed", err)
	}
}

func TestPublicationSession_ReceiveHonorsCancellation(t *testing.T) {
	// Given
	session, err := NewPublicationSession(1)
	if err != nil {
		t.Fatalf("create publication session: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// When
	_, err = session.Receive(ctx)

	// Then
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("receive error = %v, want context.Canceled", err)
	}
}
