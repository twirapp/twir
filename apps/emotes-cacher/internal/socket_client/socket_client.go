package socket_client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/goccy/go-json"
)

type WsConnection struct {
	websocket          *websocket.Conn
	socketUrl          string
	SubscriptionsCount int
	SubscriptionsLimit int
	mu                 sync.RWMutex
	ctx                context.Context
	cancel             context.CancelFunc
	onMessage          func(
		ctx context.Context,
		client *WsConnection,
		data []byte,
	)
	onReconnect func(
		ctx context.Context,
		client *WsConnection,
	)
	onConnect func(ctx context.Context, client *WsConnection)
}

type Opts struct {
	OnMessage func(
		ctx context.Context,
		client *WsConnection,
		data []byte,
	)
	OnReconnect func(
		ctx context.Context,
		client *WsConnection,
	)
	OnConnect          func(ctx context.Context, client *WsConnection)
	Url                string
	SubscriptionsLimit int
}

func connDial(ctx context.Context, url string) (*websocket.Conn, error) {
	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func New(ctx context.Context, opts Opts) (*WsConnection, error) {
	connCtx, cancel := context.WithCancel(ctx)

	wsConn := &WsConnection{
		websocket:          nil,
		SubscriptionsCount: 0,
		SubscriptionsLimit: opts.SubscriptionsLimit,
		ctx:                connCtx,
		cancel:             cancel,
		onMessage:          opts.OnMessage,
		onReconnect:        opts.OnReconnect,
		onConnect:          opts.OnConnect,
		socketUrl:          opts.Url,
		mu:                 sync.RWMutex{},
	}

	if err := wsConn.connect(); err != nil {
		cancel()
		return nil, err
	}

	go wsConn.readLoop()

	return wsConn, nil
}

func (c *WsConnection) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, err := connDial(c.ctx, c.socketUrl)
	if err != nil {
		return err
	}

	c.websocket = conn

	if c.onConnect != nil {
		c.onConnect(c.ctx, c)
	}

	return nil
}

func (c *WsConnection) reconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.websocket != nil {
		c.websocket.Close(websocket.StatusNormalClosure, "")
		c.websocket = nil
	}

	conn, err := connDial(c.ctx, c.socketUrl)
	if err != nil {
		return err
	}

	c.websocket = conn
	if c.onReconnect != nil {
		c.onReconnect(c.ctx, c)
	}

	return nil
}

func (c *WsConnection) readLoop() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		if c.websocket == nil {
			time.Sleep(time.Second)
			continue
		}

		_, msg, err := c.websocket.Read(c.ctx)
		if err != nil {
			fmt.Println("websocket connection closed, attempting to reconnect:", err)
			if reconnectErr := c.reconnect(); reconnectErr != nil {
				time.Sleep(5 * time.Second)
			}
			continue
		}

		if len(msg) == 0 {
			continue
		}

		if c.onMessage != nil {
			c.onMessage(c.ctx, c, msg)
		}
	}
}

func (c *WsConnection) Close() error {
	c.cancel()

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.websocket == nil {
		return nil
	}

	err := c.websocket.Close(websocket.StatusNormalClosure, "manual closing")
	c.websocket = nil

	return err
}

var ErrToManySubscriptions = fmt.Errorf("too many subscriptions")

func (c *WsConnection) Subscribe(ctx context.Context, msg map[string]any) error {
	if c.websocket == nil {
		return fmt.Errorf("websocket connection is nil")
	}

	if c.SubscriptionsLimit > 0 && c.SubscriptionsCount == c.SubscriptionsLimit {
		return ErrToManySubscriptions
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := c.websocket.Write(ctx, websocket.MessageText, msgBytes); err != nil {
		return fmt.Errorf("failed to send subscribe message: %w", err)
	}

	c.mu.Lock()
	c.SubscriptionsCount++
	c.mu.Unlock()

	return nil
}

func (c *WsConnection) SendMessage(ctx context.Context, msg map[string]any) error {
	if c.websocket == nil {
		return fmt.Errorf("websocket connection is nil")
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := c.websocket.Write(ctx, websocket.MessageText, msgBytes); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
