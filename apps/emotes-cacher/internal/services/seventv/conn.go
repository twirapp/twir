package seventv

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/satont/twir/apps/emotes-cacher/internal/services/seventv/messages"
)

type wsConnection struct {
	websocket          *websocket.Conn
	subscriptionsCount int
	subscriptionsLimit int
	sessionId          string
	mu                 sync.RWMutex
	ctx                context.Context
	cancel             context.CancelFunc
	onMessage          func([]byte)
}

func connDial(ctx context.Context) (*websocket.Conn, error) {
	conn, _, err := websocket.Dial(ctx, "wss://events.7tv.io/v3", nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func createConn(ctx context.Context, onMessage func([]byte)) (*wsConnection, error) {
	connCtx, cancel := context.WithCancel(ctx)

	wsConn := &wsConnection{
		websocket:          nil,
		subscriptionsCount: 0,
		subscriptionsLimit: 0,
		ctx:                connCtx,
		cancel:             cancel,
		onMessage:          onMessage,
	}

	if err := wsConn.connect(); err != nil {
		cancel()
		return nil, err
	}

	go wsConn.readLoop()

	return wsConn, nil
}

func (c *wsConnection) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, err := connDial(c.ctx)
	if err != nil {
		return err
	}

	c.websocket = conn

	// Wait for the hello message (opcode 1)
	for {
		_, msg, err := conn.Read(c.ctx)
		if err != nil {
			return err
		}

		if len(msg) == 0 {
			continue
		}

		var baseMsg messages.BaseMessage[messages.HelloMessage]
		if err := json.Unmarshal(msg, &baseMsg); err != nil {
			continue
		}

		if baseMsg.Operation == 1 { // Hello message
			c.subscriptionsLimit = int(baseMsg.Data.SubscriptionLimit)
			c.sessionId = baseMsg.Data.SessionID

			fmt.Println("connected to 7TV websocket with session ID:", c.sessionId)
			break
		}
	}

	return nil
}

func (c *wsConnection) reconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.websocket != nil {
		c.websocket.Close(websocket.StatusNormalClosure, "")
		c.websocket = nil
	}

	time.Sleep(time.Second) // Brief delay before reconnecting

	conn, err := connDial(c.ctx)
	if err != nil {
		return err
	}

	c.websocket = conn

	// Send resume message (opcode 34) if we have a session ID
	if c.sessionId != "" {
		resumeMsg := map[string]interface{}{
			"op": 34,
			"d": map[string]string{
				"session_id": c.sessionId,
			},
		}

		msgBytes, err := json.Marshal(resumeMsg)
		if err != nil {
			return err
		}

		if err := conn.Write(c.ctx, websocket.MessageText, msgBytes); err != nil {
			return err
		}

		fmt.Println("sent resume message for session:", c.sessionId)
	}

	return nil
}

func (c *wsConnection) readLoop() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		c.mu.RLock()
		conn := c.websocket
		c.mu.RUnlock()

		if conn == nil {
			time.Sleep(time.Second)
			continue
		}

		_, msg, err := conn.Read(c.ctx)
		if err != nil {
			if websocket.CloseStatus(err) != -1 {
				fmt.Println("websocket connection closed, attempting to reconnect:", err)
				if reconnectErr := c.reconnect(); reconnectErr != nil {
					fmt.Println("failed to reconnect:", reconnectErr)
					time.Sleep(5 * time.Second) // Wait before retry
				}
			}
			continue
		}

		if len(msg) == 0 {
			continue
		}

		// Parse base message to check opcode
		var baseMsg messages.BaseMessage[json.RawMessage]
		if err := json.Unmarshal(msg, &baseMsg); err != nil {
			fmt.Println("failed to parse message:", err)
			continue
		}

		switch baseMsg.Operation {
		case 1: // Hello - already handled in connect()
			continue
		case 2: // Heartbeat ACK
			fmt.Println("received heartbeat ack")
			continue
		case 0:
			if c.onMessage != nil {
				c.onMessage(msg)
			}
		default:
			fmt.Println("received unknown message type:", baseMsg.Operation, string(msg))
		}
	}
}

func (c *wsConnection) Close() error {
	c.cancel()

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.websocket == nil {
		return nil
	}

	err := c.websocket.Close(websocket.StatusNormalClosure, "")
	c.websocket = nil

	return err
}

func (c *wsConnection) subscribe(ctx context.Context, msg map[string]any) error {
	c.mu.RLock()
	conn := c.websocket
	c.mu.RUnlock()

	if conn == nil {
		return fmt.Errorf("websocket connection is nil")
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := conn.Write(ctx, websocket.MessageText, msgBytes); err != nil {
		return fmt.Errorf("failed to send subscribe message: %w", err)
	}

	c.mu.Lock()
	c.subscriptionsCount++
	c.mu.Unlock()

	fmt.Println("subscribed with message:", string(msgBytes))
	return nil
}
