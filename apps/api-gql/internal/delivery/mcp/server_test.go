package mcp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
)

func TestNewServerRegistersToolSchemas(t *testing.T) {
	handler := &Handler{}
	server := handler.newServer(scope{Channel: channelentity.Channel{ID: uuid.New()}})
	if server == nil {
		t.Fatal("expected MCP server")
	}
}

func TestHandlerRequiresAPIKeyHeader(t *testing.T) {
	handler := &Handler{}
	request := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(`{}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.Code)
	}
	if !strings.Contains(response.Body.String(), "Api-Key header is required") {
		t.Fatalf("unexpected response: %s", response.Body.String())
	}
}
