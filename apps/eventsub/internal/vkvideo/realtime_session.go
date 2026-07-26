package vkvideo

import (
	"bytes"
	"context"
	"errors"
	"sync"
)

var (
	ErrInvalidPublicationQueueCapacity = errors.New("vk video publication queue capacity must be positive")
	ErrPublicationSessionClosed        = errors.New("vk video publication session closed")
)

type PublicationSession struct {
	queue     chan []byte
	closed    chan struct{}
	closeOnce sync.Once
}

func NewPublicationSession(capacity int) (*PublicationSession, error) {
	if capacity <= 0 {
		return nil, ErrInvalidPublicationQueueCapacity
	}

	return &PublicationSession{
		queue:  make(chan []byte, capacity),
		closed: make(chan struct{}),
	}, nil
}

func (s *PublicationSession) Enqueue(publication []byte) bool {
	publicationCopy := bytes.Clone(publication)

	select {
	case <-s.closed:
		return false
	default:
	}

	select {
	case s.queue <- publicationCopy:
		return true
	case <-s.closed:
		return false
	default:
		return false
	}
}

func (s *PublicationSession) Receive(ctx context.Context) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-s.closed:
		return nil, ErrPublicationSessionClosed
	default:
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-s.closed:
		return nil, ErrPublicationSessionClosed
	case publication := <-s.queue:
		return publication, nil
	}
}

func (s *PublicationSession) Close() {
	s.closeOnce.Do(func() {
		close(s.closed)
	})
}
