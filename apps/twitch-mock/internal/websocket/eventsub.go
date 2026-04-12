package websocket

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	gws "github.com/gorilla/websocket"
	"github.com/twirapp/twir/apps/twitch-mock/internal/config"
	"github.com/twirapp/twir/apps/twitch-mock/internal/state"
)

type Server struct {
	logger   *slog.Logger
	state    *state.State
	upgrader gws.Upgrader
	mu       sync.RWMutex
	clients  map[string]*client
	mux      *http.ServeMux
}

type client struct {
	id   string
	conn *gws.Conn
	mu   sync.Mutex
	once sync.Once
}

type baseMetadata struct {
	MessageID        string `json:"message_id"`
	MessageType      string `json:"message_type"`
	MessageTimestamp string `json:"message_timestamp"`
}

type notificationMetadata struct {
	MessageID           string `json:"message_id"`
	MessageType         string `json:"message_type"`
	MessageTimestamp    string `json:"message_timestamp"`
	SubscriptionType    string `json:"subscription_type"`
	SubscriptionVersion string `json:"subscription_version"`
}

type welcomeMessage struct {
	Metadata baseMetadata `json:"metadata"`
	Payload  struct {
		Session struct {
			ID                      string  `json:"id"`
			Status                  string  `json:"status"`
			KeepaliveTimeoutSeconds int     `json:"keepalive_timeout_seconds"`
			ReconnectURL            *string `json:"reconnect_url"`
			ConnectedAt             string  `json:"connected_at"`
		} `json:"session"`
	} `json:"payload"`
}

type keepaliveMessage struct {
	Metadata baseMetadata   `json:"metadata"`
	Payload  map[string]any `json:"payload"`
}

type notificationMessage struct {
	Metadata notificationMetadata `json:"metadata"`
	Payload  struct {
		Subscription subscription   `json:"subscription"`
		Event        map[string]any `json:"event"`
	} `json:"payload"`
}

type subscription struct {
	ID        string                `json:"id"`
	Status    string                `json:"status"`
	Type      string                `json:"type"`
	Version   string                `json:"version"`
	Cost      int                   `json:"cost"`
	Condition map[string]any        `json:"condition"`
	CreatedAt string                `json:"created_at"`
	Transport subscriptionTransport `json:"transport"`
}

type subscriptionTransport struct {
	Method    string `json:"method"`
	SessionID string `json:"session_id,omitempty"`
}

func New(logger *slog.Logger, appState *state.State) *Server {
	server := &Server{
		logger: logger,
		state:  appState,
		upgrader: gws.Upgrader{
			CheckOrigin: func(*http.Request) bool { return true },
		},
		clients: map[string]*client{},
		mux:     http.NewServeMux(),
	}

	server.mux.HandleFunc("/ws", server.serveWS)

	return server
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) Broadcast(eventType string, eventData map[string]any) {
	clients := s.snapshotClients()
	for _, currentClient := range clients {
		message := s.buildNotification(currentClient.id, eventType, eventData)
		if err := currentClient.writeJSON(message); err != nil {
			s.logger.Warn("failed to broadcast websocket notification", slog.String("session_id", currentClient.id), slog.Any("error", err))
			s.removeClient(currentClient.id)
			currentClient.close()
		}
	}
}

func (s *Server) serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Warn("failed to upgrade websocket connection", slog.Any("error", err))
		return
	}

	currentClient := &client{
		id:   uuid.NewString(),
		conn: conn,
	}

	s.addClient(currentClient)

	if err := currentClient.writeJSON(s.buildWelcome(currentClient.id)); err != nil {
		s.logger.Warn("failed to send websocket welcome", slog.Any("error", err))
		s.removeClient(currentClient.id)
		currentClient.close()
		return
	}

	go s.keepaliveLoop(currentClient)
	s.readLoop(currentClient)
}

func (s *Server) keepaliveLoop(currentClient *client) {
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := currentClient.writeJSON(buildKeepalive()); err != nil {
			s.logger.Info("closing websocket client after keepalive failure", slog.String("session_id", currentClient.id), slog.Any("error", err))
			s.removeClient(currentClient.id)
			currentClient.close()
			return
		}
	}
}

func (s *Server) readLoop(currentClient *client) {
	defer func() {
		s.removeClient(currentClient.id)
		currentClient.close()
	}()

	for {
		if _, _, err := currentClient.conn.ReadMessage(); err != nil {
			return
		}
	}
}

func (s *Server) buildWelcome(sessionID string) welcomeMessage {
	now := time.Now().UTC().Format(time.RFC3339)
	message := welcomeMessage{
		Metadata: baseMetadata{
			MessageID:        uuid.NewString(),
			MessageType:      "session_welcome",
			MessageTimestamp: now,
		},
	}

	message.Payload.Session.ID = sessionID
	message.Payload.Session.Status = "connected"
	message.Payload.Session.KeepaliveTimeoutSeconds = 30
	message.Payload.Session.ConnectedAt = now

	return message
}

func buildKeepalive() keepaliveMessage {
	return keepaliveMessage{
		Metadata: baseMetadata{
			MessageID:        uuid.NewString(),
			MessageType:      "session_keepalive",
			MessageTimestamp: time.Now().UTC().Format(time.RFC3339),
		},
		Payload: map[string]any{},
	}
}

func (s *Server) buildNotification(sessionID string, eventType string, eventData map[string]any) notificationMessage {
	subscriptionData, ok := s.state.FindSubscriptionByType(eventType)
	if !ok {
		subscriptionData = state.Subscription{
			ID:        uuid.NewString(),
			Status:    "enabled",
			Type:      eventType,
			Version:   defaultVersion(eventType),
			Cost:      0,
			Condition: map[string]any{"broadcaster_user_id": config.MockBroadcasterID, "moderator_user_id": config.MockBotID},
			CreatedAt: time.Now().UTC(),
			Transport: state.SubscriptionTransport{Method: "websocket"},
		}
	}

	message := notificationMessage{
		Metadata: notificationMetadata{
			MessageID:           uuid.NewString(),
			MessageType:         "notification",
			MessageTimestamp:    time.Now().UTC().Format(time.RFC3339),
			SubscriptionType:    eventType,
			SubscriptionVersion: subscriptionData.Version,
		},
	}
	message.Payload.Event = cloneMap(eventData)

	message.Payload.Subscription = subscription{
		ID:        subscriptionData.ID,
		Status:    subscriptionData.Status,
		Type:      subscriptionData.Type,
		Version:   subscriptionData.Version,
		Cost:      subscriptionData.Cost,
		Condition: cloneMap(subscriptionData.Condition),
		CreatedAt: subscriptionData.CreatedAt.UTC().Format(time.RFC3339),
		Transport: subscriptionTransport{
			Method:    "websocket",
			SessionID: sessionID,
		},
	}

	return message
}

func (s *Server) addClient(currentClient *client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[currentClient.id] = currentClient
}

func (s *Server) removeClient(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.clients, sessionID)
}

func (s *Server) snapshotClients() []*client {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*client, 0, len(s.clients))
	for _, currentClient := range s.clients {
		result = append(result, currentClient)
	}

	return result
}

func (c *client) writeJSON(payload any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.WriteJSON(payload)
}

func (c *client) close() {
	c.once.Do(func() {
		_ = c.conn.Close()
	})
}

func defaultVersion(eventType string) string {
	switch eventType {
	case "channel.follow":
		return "2"
	default:
		return "1"
	}
}

func cloneMap(input map[string]any) map[string]any {
	if input == nil {
		return map[string]any{}
	}

	output := make(map[string]any, len(input))
	for key, value := range input {
		output[key] = value
	}

	return output
}
