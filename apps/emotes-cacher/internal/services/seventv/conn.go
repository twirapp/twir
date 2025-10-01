package seventv

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/tmaxmax/go-sse"
	"github.com/twirapp/twir/apps/emotes-cacher/internal/services/seventv/messages"
	"github.com/twirapp/twir/apps/emotes-cacher/internal/services/seventv/operations"
)

type conn struct {
	restartCh chan struct{}
	stopCh    chan struct{}

	subscriptionsMu sync.Mutex
	subscriptions   []connSubscription
	maxCapacity     int
	onMessage       onConnMessage
}

type connSubscription struct {
	subType    string
	conditions map[string]string
}

type onConnMessage func(
	ctx context.Context,
	client *conn,
	msg []byte,
)

func newConn(
	onMessage onConnMessage,
	maxCapacity int,
) *conn {
	c := &conn{
		subscriptionsMu: sync.Mutex{},
		subscriptions:   make([]connSubscription, 0),
		maxCapacity:     maxCapacity,
		onMessage:       onMessage,
		stopCh:          make(chan struct{}),
		restartCh:       make(chan struct{}),
	}

	go c.start()

	return c
}

func (c *conn) createConnUrl() string {
	baseUrl, _ := url.Parse("https://events.7tv.io/v3")

	for i, sub := range c.subscriptions {
		if i == 0 {
			baseUrl.Path += "@"
		} else {
			baseUrl.Path += ","
		}

		baseUrl.Path += sub.subType + "<"
		j := 0
		for k, v := range sub.conditions {
			if j != 0 {
				baseUrl.Path += ";"
			}
			baseUrl.Path += k + "=" + v
			j++
		}
		baseUrl.Path += ">"
	}

	return baseUrl.String()
}

func (c *conn) start() error {
	reqCtx, reqCtxCancel := context.WithCancel(context.Background())

	go func() {
		var mu sync.Mutex
		for {
			select {
			case <-c.stopCh:
				mu.Lock()
				reqCtxCancel()
				mu.Unlock()
				return
			case <-c.restartCh:
				mu.Lock()
				reqCtxCancel()
				reqCtx, reqCtxCancel = context.WithCancel(context.Background())
				mu.Unlock()
				continue
			}
		}
	}()

	go func() {
		for {
			connUrl := c.createConnUrl()
			r, _ := http.NewRequestWithContext(reqCtx, http.MethodGet, connUrl, nil)
			r.Header.Set("Accept", "text/event-stream")
			r.Header.Set("Cache-Control", "no-cache")
			r.Header.Set("Connection", "keep-alive")

			client := sse.Client{
				Backoff: sse.Backoff{
					MaxRetries: -1,
				},
			}
			clientConn := client.NewConnection(r)
			clientConn.SubscribeToAll(
				func(event sse.Event) {
					op, ok := operations.OpFromString(event.Type)
					if !ok {
						log.Println("[7TV] Unknown operation:", event.Type)
						return
					}

					msg := messages.BaseMessage[json.RawMessage]{
						Data:      []byte(event.Data),
						Operation: int(op),
					}

					msgJson, err := json.Marshal(msg)
					if err != nil {
						log.Println("[7TV] JSON marshal error:", err)
						return
					}

					c.onMessage(reqCtx, c, msgJson)
				},
			)

			if err := clientConn.Connect(); err != nil && !errors.Is(err, context.Canceled) {
				log.Println("[7TV] Connection error:", err)
			}

			time.Sleep(500 * time.Millisecond)
		}
	}()

	return nil
}

func (c *conn) Stop() error {
	c.stopCh <- struct{}{}

	return nil
}

func (c *conn) Restart() error {
	c.restartCh <- struct{}{}

	return nil
}

func (c *conn) subscribe(subs ...connSubscription) error {
	if len(c.subscriptions)+len(subs) > c.maxCapacity {
		return fmt.Errorf("max capacity reached")
	}

	c.subscriptionsMu.Lock()
	c.subscriptions = append(
		c.subscriptions,
		subs...,
	)
	c.subscriptionsMu.Unlock()

	return c.Restart()
}
