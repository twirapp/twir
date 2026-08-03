package mcp

import (
	"context"
	"testing"

	"github.com/google/uuid"
	modelsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	oauthentity "github.com/twirapp/twir/libs/entities/mcp_oauth"
)

const registeredToolCount = 85

func TestToolAccessScopesAllowsOnlyListAndGetToolsWithReadScope(t *testing.T) {
	tests := []struct {
		name  string
		tool  string
		allow bool
	}{
		{name: "allows list tool", tool: "list_commands", allow: true},
		{name: "allows get tool", tool: "get_command", allow: true},
		{name: "rejects sensitive get tool", tool: "get_secret", allow: false},
		{name: "rejects write-required list tool", tool: "list_overlays", allow: false},
		{name: "rejects write-required get tool", tool: "get_overlay", allow: false},
		{name: "rejects create tool", tool: "create_command", allow: false},
		{name: "rejects evaluation tool", tool: "evaluate_variable", allow: false},
		{name: "rejects similarly named tool", tool: "listing_commands", allow: false},
	}
	readScopes := toolAccessScopes{toolAccessScopeRead: {}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			scopes := readScopes

			// When
			got := scopes.allowsTool(test.tool)

			// Then
			if got != test.allow {
				t.Fatalf("allowsTool(%q) = %t, want %t", test.tool, got, test.allow)
			}
		})
	}
}

func TestClassifyToolAccessMarksSensitiveToolsBeforeReadPrefix(t *testing.T) {
	tests := []struct {
		name string
		tool string
		want toolAccessClassification
	}{
		{name: "sensitive get tool", tool: "get_secret", want: toolAccessClassificationSensitive},
		{name: "write-required list tool", tool: "list_overlays", want: toolAccessClassificationWriteRequired},
		{name: "write-required get tool", tool: "get_overlay", want: toolAccessClassificationWriteRequired},
		{name: "read list tool", tool: "list_commands", want: toolAccessClassificationRead},
		{name: "read get tool", tool: "get_command", want: toolAccessClassificationRead},
		{name: "other tool", tool: "create_command", want: toolAccessClassificationOther},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := classifyToolAccess(test.tool); got != test.want {
				t.Fatalf("classifyToolAccess(%q) = %q, want %q", test.tool, got, test.want)
			}
		})
	}
}

func TestToolAccessScopesAllowsAllToolsWithWriteScope(t *testing.T) {
	// Given
	scopes := toolAccessScopes{toolAccessScopeWrite: {}}

	// When / Then
	for _, name := range []string{"list_commands", "get_command", "create_command", "evaluate_variable", "list_overlays", "get_overlay"} {
		if !scopes.allowsTool(name) {
			t.Fatalf("write scope did not allow %q", name)
		}
	}
	if !scopes.allowsTool("get_secret") {
		t.Fatal("write scope did not allow get_secret")
	}
}

func TestNewServerListsOnlyReadToolsWithReadScope(t *testing.T) {
	// Given
	tools := listServerTools(t, toolAccessScopes{toolAccessScopeRead: {}})

	// When / Then
	for name := range tools {
		if classifyToolAccess(name) != toolAccessClassificationRead {
			t.Fatalf("read scope exposed write tool %q", name)
		}
	}
	for _, name := range []string{"list_commands", "get_command", "get_integration_status"} {
		if _, ok := tools[name]; !ok {
			t.Fatalf("read scope did not expose %q", name)
		}
	}
	for _, name := range []string{"get_secret", "list_overlays", "get_overlay"} {
		if _, ok := tools[name]; ok {
			t.Fatalf("read scope exposed %q", name)
		}
	}
	for _, name := range []string{"create_command", "evaluate_variable"} {
		if _, ok := tools[name]; ok {
			t.Fatalf("read scope exposed %q", name)
		}
	}
}

func TestNewServerListsCompleteSurfaceWithWriteScope(t *testing.T) {
	// Given
	readTools := listServerTools(t, toolAccessScopes{toolAccessScopeRead: {}})
	fullScopes, ok := toolAccessScopesFromOAuthScopes([]oauthentity.Scope{oauthentity.ScopeRead, oauthentity.ScopeWrite})
	if !ok {
		t.Fatal("expected write OAuth grant to map to tool access scopes")
	}
	fullTools := listServerTools(t, fullScopes)

	// When / Then
	if len(fullTools) != registeredToolCount {
		t.Fatalf("full scope registered %d tools, want %d", len(fullTools), registeredToolCount)
	}
	for name := range readTools {
		if _, ok := fullTools[name]; !ok {
			t.Fatalf("full scope omitted read tool %q", name)
		}
	}
	for _, name := range []string{"create_command", "evaluate_variable"} {
		if _, ok := fullTools[name]; !ok {
			t.Fatalf("full scope omitted %q", name)
		}
	}
	for _, name := range []string{"get_secret", "list_overlays", "get_overlay"} {
		if _, ok := fullTools[name]; !ok {
			t.Fatalf("full scope omitted %q", name)
		}
	}
}

func TestToolAccessScopesFromOAuthScopesRejectsInvalidGrants(t *testing.T) {
	tests := []struct {
		name   string
		scopes []oauthentity.Scope
	}{
		{name: "no scopes"},
		{name: "write without read", scopes: []oauthentity.Scope{oauthentity.ScopeWrite}},
		{name: "unknown scope", scopes: []oauthentity.Scope{oauthentity.ScopeRead, "admin"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, ok := toolAccessScopesFromOAuthScopes(test.scopes)
			if ok {
				t.Fatalf("scopes %v unexpectedly mapped to tool access", test.scopes)
			}
		})
	}
}

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
