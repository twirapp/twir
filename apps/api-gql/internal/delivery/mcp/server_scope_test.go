package mcp

import (
	"net/http"
	"strings"
	"testing"

	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
)

func TestHandlerLegacyReadGrantExposesMappedReadTools(t *testing.T) {
	verifier := &mcpAccessVerifier{grant: testAuthorizedGrant(entity.ScopeRead)}
	body := mcpToolsList(t, newTestMCPHandler(verifier))

	for name, permission := range toolPermissions {
		toolListed := strings.Contains(body, `"`+name+`"`)
		if permission.Action == entity.ScopeActionRead && !toolListed {
			t.Errorf("read tool %q was not listed", name)
		}
		if permission.Action == entity.ScopeActionEdit && toolListed {
			t.Errorf("edit tool %q was listed for read grant", name)
		}
	}
}

func TestHandlerLegacyWriteGrantExposesCompleteToolSurface(t *testing.T) {
	verifier := &mcpAccessVerifier{grant: testAuthorizedGrant(entity.ScopeWrite)}
	body := mcpToolsList(t, newTestMCPHandler(verifier))

	if len(toolPermissions) != 85 {
		t.Fatalf("mapped tool count = %d, want 85", len(toolPermissions))
	}
	for name := range toolPermissions {
		if !strings.Contains(body, `"`+name+`"`) {
			t.Errorf("tool %q was not listed for legacy write grant", name)
		}
	}
}

func TestHandlerCommandScopesExposeOnlyCommandTools(t *testing.T) {
	tests := []struct {
		name      string
		scope     entity.Scope
		wanted    []string
		forbidden []string
	}{
		{name: "read", scope: "commands:read", wanted: []string{"list_command_roles", "list_command_groups", "list_commands", "get_command"}, forbidden: []string{"create_command", "update_command", "delete_command", "list_timers", "list_secrets"}},
		{name: "edit", scope: "commands:edit", wanted: []string{"list_command_roles", "list_command_groups", "list_commands", "get_command", "create_command", "update_command", "delete_command"}, forbidden: []string{"list_timers", "list_secrets"}},
	}

	assertToolSurface := func(t *testing.T, test struct {
		name      string
		scope     entity.Scope
		wanted    []string
		forbidden []string
	}) {
		t.Helper()
		verifier := &mcpAccessVerifier{grant: testAuthorizedGrant(test.scope)}
		body := mcpToolsList(t, newTestMCPHandler(verifier))
		for _, name := range test.wanted {
			if !strings.Contains(body, `"`+name+`"`) {
				t.Errorf("tool %q was not listed", name)
			}
		}
		for _, name := range test.forbidden {
			if strings.Contains(body, `"`+name+`"`) {
				t.Errorf("tool %q was listed", name)
			}
		}
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) { assertToolSurface(t, test) })
	}
}

func TestHandlerSecretScopesExposeExpectedTools(t *testing.T) {
	tests := []struct {
		name      string
		scope     entity.Scope
		wanted    []string
		forbidden []string
	}{
		{name: "read", scope: "secrets:read", wanted: []string{"list_secrets"}, forbidden: []string{"get_secret", "create_secret", "update_secret", "delete_secret"}},
		{name: "edit", scope: "secrets:edit", wanted: []string{"list_secrets", "get_secret", "create_secret", "update_secret", "delete_secret"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			verifier := &mcpAccessVerifier{grant: testAuthorizedGrant(test.scope)}
			body := mcpToolsList(t, newTestMCPHandler(verifier))
			for _, name := range test.wanted {
				if !strings.Contains(body, `"`+name+`"`) {
					t.Errorf("tool %q was not listed", name)
				}
			}
			for _, name := range test.forbidden {
				if strings.Contains(body, `"`+name+`"`) {
					t.Errorf("tool %q was listed", name)
				}
			}
		})
	}
}

func TestHandlerOverlayScopesExposeExpectedTools(t *testing.T) {
	tests := []struct {
		name      string
		scope     entity.Scope
		wanted    []string
		forbidden []string
	}{
		{name: "read", scope: "overlays:read", forbidden: []string{"list_overlays", "get_overlay"}},
		{name: "edit", scope: "overlays:edit", wanted: []string{"list_overlays", "get_overlay"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			verifier := &mcpAccessVerifier{grant: testAuthorizedGrant(test.scope)}
			body := mcpToolsList(t, newTestMCPHandler(verifier))
			for _, name := range test.wanted {
				if !strings.Contains(body, `"`+name+`"`) {
					t.Errorf("tool %q was not listed", name)
				}
			}
			for _, name := range test.forbidden {
				if strings.Contains(body, `"`+name+`"`) {
					t.Errorf("tool %q was listed", name)
				}
			}
		})
	}
}

func mcpToolsList(t *testing.T, handler *Handler) string {
	t.Helper()
	response := serveMCPRequest(handler, `{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}`)
	if response.Code != http.StatusOK {
		t.Fatalf("tools/list response = %d %s", response.Code, response.Body.String())
	}
	return response.Body.String()
}
