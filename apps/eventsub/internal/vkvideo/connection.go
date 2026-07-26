package vkvideo

import (
	"context"
	"sync"
)

type realtimeConnection interface {
	Connect() error
	Receive(context.Context) ([]byte, error)
	Close()
}

type realtimeConnectionFactory func(RealtimeClientConfig) (realtimeConnection, error)

type activeBinding struct {
	lease      *Lease
	connection *ownedConnection
}

type ownedConnection struct {
	mu     sync.Mutex
	closed bool
	conn   realtimeConnection
}

func (c *ownedConnection) Set(connection realtimeConnection) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		connection.Close()
		return
	}
	c.conn = connection
}

func (c *ownedConnection) Receive(ctx context.Context) ([]byte, error) {
	c.mu.Lock()
	connection := c.conn
	c.mu.Unlock()
	if connection == nil {
		return nil, ErrPublicationSessionClosed
	}
	return connection.Receive(ctx)
}

func (c *ownedConnection) Close() {
	c.mu.Lock()
	c.closed = true
	connection := c.conn
	c.mu.Unlock()
	if connection != nil {
		connection.Close()
	}
}

type centrifugoConnection struct {
	client *RealtimeClient
}

func newCentrifugoConnection(config RealtimeClientConfig) (realtimeConnection, error) {
	client, err := NewRealtimeClient(config)
	if err != nil {
		return nil, err
	}
	return centrifugoConnection{client: client}, nil
}

func (c centrifugoConnection) Connect() error {
	return c.client.Connect()
}

func (c centrifugoConnection) Receive(ctx context.Context) ([]byte, error) {
	return c.client.Session().Receive(ctx)
}

func (c centrifugoConnection) Close() {
	c.client.Close()
}
