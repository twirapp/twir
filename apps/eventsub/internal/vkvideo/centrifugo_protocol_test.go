package vkvideo

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/centrifugal/protocol"
	"github.com/gorilla/websocket"
)

const protocolCommandTimeout = 5 * time.Second

type centrifugoProtocolServer struct {
	t *testing.T

	server   *httptest.Server
	endpoint string
	channel  string

	connectionsMu sync.Mutex
	connection    *websocket.Conn
	writesMu      sync.Mutex
	lifecycleMu   sync.Mutex
	closing       bool
	connections   sync.WaitGroup
	closeOnce     sync.Once

	connectCommands   chan *protocol.Command
	subscribeCommands chan *protocol.Command
	connectAck        chan struct{}
	subscribeAck      chan struct{}
	connectAckOnce    sync.Once
	subscribeAckOnce  sync.Once
	closed            chan struct{}
	errors            chan error
}

func newCentrifugoProtocolServer(t *testing.T, channel string) *centrifugoProtocolServer {
	t.Helper()

	server := &centrifugoProtocolServer{
		t:                 t,
		channel:           channel,
		connectCommands:   make(chan *protocol.Command, 1),
		subscribeCommands: make(chan *protocol.Command, 1),
		connectAck:        make(chan struct{}),
		subscribeAck:      make(chan struct{}),
		closed:            make(chan struct{}),
		errors:            make(chan error, 1),
	}
	server.server = httptest.NewServer(http.HandlerFunc(server.handleConnection))
	server.endpoint = "ws" + strings.TrimPrefix(server.server.URL, "http") + "/connection/websocket?cf_protocol_version=v2"
	t.Cleanup(server.Close)

	return server
}

func (s *centrifugoProtocolServer) handleConnection(writer http.ResponseWriter, request *http.Request) {
	s.lifecycleMu.Lock()
	if s.closing {
		s.lifecycleMu.Unlock()
		return
	}
	s.connections.Add(1)
	s.lifecycleMu.Unlock()
	defer s.connections.Done()

	connection, err := (&websocket.Upgrader{}).Upgrade(writer, request, nil)
	if err != nil {
		s.reportError(errors.New("protocol websocket upgrade failed"))
		return
	}

	s.connectionsMu.Lock()
	s.connection = connection
	s.connectionsMu.Unlock()

	select {
	case <-s.closed:
		if err := connection.Close(); err != nil {
			s.t.Errorf("close test websocket: %v", err)
		}
		return
	default:
	}

	s.readCommands(connection)
}

func (s *centrifugoProtocolServer) readCommands(connection *websocket.Conn) {
	for {
		messageType, data, err := connection.ReadMessage()
		if err != nil {
			return
		}
		if messageType != websocket.TextMessage {
			s.reportError(errors.New("protocol command used a non-text frame"))
			return
		}

		decoder := protocol.NewJSONCommandDecoder(data)
		for {
			command, decodeErr := decoder.Decode()
			if command != nil {
				s.handleCommand(connection, command)
			}
			if errors.Is(decodeErr, io.EOF) {
				break
			}
			if decodeErr != nil {
				s.reportError(errors.New("protocol command decode failed"))
				return
			}
		}
	}
}

func (s *centrifugoProtocolServer) handleCommand(connection *websocket.Conn, command *protocol.Command) {
	switch {
	case command.Connect != nil:
		if !s.sendCommand(s.connectCommands, command) {
			return
		}
		select {
		case <-s.connectAck:
			if err := s.sendReply(connection, &protocol.Reply{
				Id:      command.Id,
				Connect: &protocol.ConnectResult{Client: "test-client", Version: "test-version"},
			}); err != nil {
				s.reportError(err)
			}
		case <-s.closed:
		}
	case command.Subscribe != nil:
		if !s.sendCommand(s.subscribeCommands, command) {
			return
		}
		select {
		case <-s.subscribeAck:
			if err := s.sendReply(connection, &protocol.Reply{
				Id:        command.Id,
				Subscribe: &protocol.SubscribeResult{},
			}); err != nil {
				s.reportError(err)
			}
		case <-s.closed:
		}
	default:
		s.reportError(errors.New("unexpected protocol command"))
	}
}

func (s *centrifugoProtocolServer) sendCommand(commands chan<- *protocol.Command, command *protocol.Command) bool {
	select {
	case commands <- command:
		return true
	case <-s.closed:
		return false
	}
}

func (s *centrifugoProtocolServer) sendReply(connection *websocket.Conn, reply *protocol.Reply) error {
	data, err := protocol.NewJSONReplyEncoder().Encode(reply)
	if err != nil {
		return errors.New("protocol reply encode failed")
	}

	s.writesMu.Lock()
	defer s.writesMu.Unlock()
	if err := connection.WriteMessage(websocket.TextMessage, data); err != nil {
		return errors.New("protocol reply write failed")
	}
	return nil
}

func (s *centrifugoProtocolServer) AwaitConnect(ctx context.Context) (*protocol.ConnectRequest, error) {
	command, err := s.awaitCommand(ctx, s.connectCommands)
	if err != nil {
		return nil, err
	}
	return command.Connect, nil
}

func (s *centrifugoProtocolServer) AwaitSubscribe(ctx context.Context) (*protocol.SubscribeRequest, error) {
	command, err := s.awaitCommand(ctx, s.subscribeCommands)
	if err != nil {
		return nil, err
	}
	return command.Subscribe, nil
}

func (s *centrifugoProtocolServer) awaitCommand(ctx context.Context, commands <-chan *protocol.Command) (*protocol.Command, error) {
	select {
	case command := <-commands:
		return command, nil
	case err := <-s.errors:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *centrifugoProtocolServer) ReleaseConnectAcknowledgement() {
	s.connectAckOnce.Do(func() { close(s.connectAck) })
}

func (s *centrifugoProtocolServer) ReleaseSubscribeAcknowledgement() {
	s.subscribeAckOnce.Do(func() { close(s.subscribeAck) })
}

func (s *centrifugoProtocolServer) RequireNoSubscribe(t *testing.T) {
	t.Helper()

	timer := time.NewTimer(150 * time.Millisecond)
	defer timer.Stop()
	select {
	case <-s.subscribeCommands:
		t.Fatal("subscription command reached the protocol server")
	case err := <-s.errors:
		t.Fatalf("protocol server error: %v", err)
	case <-timer.C:
	}
}

func (s *centrifugoProtocolServer) SendPublication(data []byte) error {
	s.connectionsMu.Lock()
	connection := s.connection
	s.connectionsMu.Unlock()
	if connection == nil {
		return errors.New("protocol websocket is not connected")
	}

	return s.sendReply(connection, &protocol.Reply{Push: &protocol.Push{
		Channel: s.channel,
		Pub:     &protocol.Publication{Data: data},
	}})
}

func (s *centrifugoProtocolServer) reportError(err error) {
	select {
	case <-s.closed:
		return
	default:
	}
	select {
	case s.errors <- err:
	default:
	}
}

func (s *centrifugoProtocolServer) Close() {
	s.t.Helper()
	s.closeOnce.Do(func() {
		s.lifecycleMu.Lock()
		s.closing = true
		s.lifecycleMu.Unlock()
		close(s.closed)

		s.connectionsMu.Lock()
		connection := s.connection
		s.connectionsMu.Unlock()
		if connection != nil {
			if err := connection.Close(); err != nil {
				s.t.Errorf("close test websocket: %v", err)
			}
		}
		s.server.Close()
		s.connections.Wait()
	})
}
