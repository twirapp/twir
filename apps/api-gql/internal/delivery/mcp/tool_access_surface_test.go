package mcp

import (
	"context"
	"testing"

	"github.com/google/uuid"
	modelsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
)

func listServerTools(t *testing.T, scopes toolAccessScopes) map[string]struct{} {
	t.Helper()

	// Given
	ctx := context.Background()
	server := (&Handler{}).newServer(scope{
		Channel:      channelentity.Channel{ID: uuid.New()},
		AccessScopes: scopes,
	})
	client := modelsdk.NewClient(&modelsdk.Implementation{Name: "test-client", Version: "1.0.0"}, nil)
	serverTransport, clientTransport := modelsdk.NewInMemoryTransports()
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatalf("connect MCP server: %v", err)
	}
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("connect MCP client: %v", err)
	}
	t.Cleanup(func() {
		if err := clientSession.Close(); err != nil {
			t.Errorf("close MCP client: %v", err)
		}
		if err := serverSession.Wait(); err != nil {
			t.Errorf("wait for MCP server: %v", err)
		}
	})

	// When
	result, err := clientSession.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list MCP tools: %v", err)
	}

	// Then
	tools := make(map[string]struct{}, len(result.Tools))
	for _, tool := range result.Tools {
		tools[tool.Name] = struct{}{}
	}
	return tools
}
