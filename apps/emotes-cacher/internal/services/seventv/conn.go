package seventv

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coder/websocket"
	"github.com/satont/twir/apps/emotes-cacher/internal/services/seventv/messages"
)

type wsConnection struct {
	websocket     *websocket.Conn
	channels      []string
	channelsLimit int
	sessionId     string
}

func connDial(ctx context.Context) (*websocket.Conn, error) {
	conn, _, err := websocket.Dial(ctx, "wss://events.7tv.io/v3", nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func createConn(ctx context.Context, onMessage func([]byte)) (*wsConnection, error) {
	wsConn := wsConnection{
		websocket:     nil,
		channels:      []string{},
		channelsLimit: 0,
	}

	conn, err := connDial(ctx)
	if err != nil {
		return nil, err
	}

	wsConn.websocket = conn

	// wait for the hello message
	for {
		_, msg, err := conn.Read(ctx)
		if err != nil {
			return nil, err
		}

		if len(msg) == 0 {
			continue
		}

		var baseMsg messages.BaseMessage[messages.HelloMessage]
		if err := json.Unmarshal(msg, &baseMsg); err != nil {
			return nil, err
		}

		if baseMsg.Operation == 1 {
			wsConn.channelsLimit = int(baseMsg.Data.SubscriptionLimit)
			wsConn.sessionId = baseMsg.Data.SessionID

			fmt.Println("connected to 7TV websocket with session ID:", wsConn.sessionId)

			break
		}
	}

	go func() {
		time.Sleep(2 * time.Second)
		wsConn.Close()
	}()

	go func() {
		for {
			_, msg, err := conn.Read(ctx)
			if err != nil {
				if websocket.CloseStatus(err) != -1 {
					time.Sleep(time.Second)
					newConn, err := connDial(ctx)
					if err != nil {
						fmt.Println("failed to reconnect:", err)
						return
					}

					wsConn.websocket = newConn
					conn = newConn

					if err := conn.Write(
						ctx,
						websocket.MessageText,
						[]byte(`{"op": 34, "d": {"session_id": "`+wsConn.sessionId+`"}}`),
					); err != nil {
						fmt.Println("failed to resume session:", err)
					}

					fmt.Println("reconnected to 7TV websocket with session ID:", wsConn.sessionId)
				}
				continue
			}
			if len(msg) == 0 {
				continue
			}

			fmt.Println("received message:", string(msg))

			onMessage(msg)
		}
	}()

	return &wsConn, nil
}

func (c *wsConnection) Close() error {
	if c.websocket == nil {
		return nil
	}

	err := c.websocket.Close(websocket.StatusNormalClosure, "")
	c.websocket = nil

	return err
}

func (c *wsConnection) addChannel(ctx context.Context, channelId string) {
	if c.websocket == nil {
		return
	}

	// wsjson.Write(
	// 	ctx,
	// 	c.websocket,
	// 	map[string]any{
	// 		"op":   "subscribe",
	// 		"data": map[string]any{"channels": []string{channelId}},
	// 	},
	// )

	c.channels = append(c.channels, channelId)
}
