package mcp

import (
	"context"
	"strings"
	"testing"

	modelsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	oauthentity "github.com/twirapp/twir/libs/entities/mcp_oauth"
)

const registeredToolCount = 85

var expectedToolPermissions = []struct {
	group oauthentity.ScopeGroup
	read  string
	edit  string
}{
	{oauthentity.ScopeGroupCommands, "list_command_roles list_command_groups list_commands get_command", "create_command update_command delete_command"},
	{oauthentity.ScopeGroupTimers, "list_timers", "create_timer update_timer delete_timer"},
	{oauthentity.ScopeGroupFiles, "list_files", "upload_file"},
	{oauthentity.ScopeGroupGames, "list_games", ""},
	{oauthentity.ScopeGroupSongRequests, "list_song_requests get_current_song", "skip_song manage_queue"},
	{oauthentity.ScopeGroupModeration, "get_moderation_settings list_mod_chat_wall", "update_moderation"},
	{oauthentity.ScopeGroupOverlays, "", "list_overlays get_overlay"},
	{oauthentity.ScopeGroupIntegrations, "get_integration_status", "toggle_integration"},
	{oauthentity.ScopeGroupEvents, "list_events list_twir_events", ""},
	{oauthentity.ScopeGroupRewards, "list_rewards", "manage_rewards"},
	{oauthentity.ScopeGroupGiveaways, "list_giveaways", "create_giveaway"},
	{oauthentity.ScopeGroupGreetings, "list_greetings", "update_greeting"},
	{oauthentity.ScopeGroupNotifications, "list_notifications get_notification", ""},
	{oauthentity.ScopeGroupAlerts, "list_alerts", "manage_alerts"},
	{oauthentity.ScopeGroupSecrets, "list_secrets", "get_secret create_secret update_secret delete_secret"},
	{oauthentity.ScopeGroupStorage, "list_storage_files get_storage_file get_storage_usage", "upload_storage_file delete_storage_file"},
	{oauthentity.ScopeGroupPastes, "list_pastes get_paste get_paste_raw", "create_paste update_paste delete_paste"},
	{oauthentity.ScopeGroupShortURLs, "get_short_url list_short_urls list_short_links get_short_url_stats list_banned_user_agents", "shorten_url update_short_url delete_short_url ban_user_agent unban_user_agent"},
	{oauthentity.ScopeGroupDashboard, "get_stats get_bot_settings list_community_users get_user_info list_scheduled_vips", "update_bot_settings"},
	{oauthentity.ScopeGroupVariables, "list_builtin_variables list_variables get_variable", "create_variable set_variable update_variable delete_variable evaluate_variable"},
	{oauthentity.ScopeGroupQuotes, "list_quotes", "create_quote update_quote delete_quote"},
	{oauthentity.ScopeGroupKeywords, "list_keywords", "create_keyword update_keyword delete_keyword"},
}

func TestToolPermissionsMatchRegisteredToolsBidirectionally(t *testing.T) {
	// Given
	wantPermissions := expectedPermissionMap()
	fullScopes, ok := toolAccessScopesFromOAuthScopes(oauthentity.AllScopes())
	if !ok {
		t.Fatal("expected all canonical OAuth scopes to map")
	}

	// When
	tools := listServerTools(t, fullScopes)

	// Then
	if len(wantPermissions) != registeredToolCount {
		t.Fatalf("expected permission specification to contain %d tools, got %d", registeredToolCount, len(wantPermissions))
	}
	if len(toolPermissions) != registeredToolCount {
		t.Fatalf("expected permission map to contain %d tools, got %d", registeredToolCount, len(toolPermissions))
	}
	if len(tools) != registeredToolCount {
		t.Fatalf("full scope registered %d tools, want %d", len(tools), registeredToolCount)
	}
	for name, want := range wantPermissions {
		got, mapped := toolPermissions[name]
		if !mapped {
			t.Fatalf("tool %q is missing from the permission map", name)
		}
		if got != want {
			t.Fatalf("permission for %q = %#v, want %#v", name, got, want)
		}
	}
	for name := range toolPermissions {
		if _, registered := tools[name]; !registered {
			t.Fatalf("permission map contains stale tool %q", name)
		}
	}
	for name := range tools {
		if _, mapped := toolPermissions[name]; !mapped {
			t.Fatalf("registered tool %q is missing from the permission map", name)
		}
	}
}

func TestToolAccessScopesEnforceReadEditAndGroupIsolation(t *testing.T) {
	for _, group := range oauthentity.AllScopeGroups() {
		for _, action := range []oauthentity.ScopeAction{oauthentity.ScopeActionRead, oauthentity.ScopeActionEdit} {
			name := string(group.Group) + ":" + string(action)
			scopes, ok := toolAccessScopesFromOAuthScopes([]oauthentity.Scope{oauthentity.Scope(name)})
			if !ok {
				t.Fatalf("scope %q did not map", name)
			}

			t.Run(name, func(t *testing.T) {
				for toolName, permission := range toolPermissions {
					want := permission.Group == group.Group && (action == oauthentity.ScopeActionEdit || permission.Action == oauthentity.ScopeActionRead)
					if got := scopes.allowsTool(toolName); got != want {
						t.Errorf("allowsTool(%q) with %s = %t, want %t", toolName, name, got, want)
					}
				}
			})
		}
	}
}

func TestLegacyScopesExposeExpectedToolSets(t *testing.T) {
	tests := []struct {
		name  string
		scope oauthentity.Scope
		want  func(string) bool
	}{
		{name: "legacy read", scope: oauthentity.ScopeRead, want: func(permission string) bool { return permission == "read" }},
		{name: "legacy write", scope: oauthentity.ScopeWrite, want: func(string) bool { return true }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			scopes, ok := toolAccessScopesFromOAuthScopes([]oauthentity.Scope{test.scope})
			if !ok {
				t.Fatalf("legacy scope %q did not map", test.scope)
			}

			// When
			tools := listServerTools(t, scopes)

			// Then
			for name, permission := range toolPermissions {
				want := test.want(string(permission.Action))
				_, got := tools[name]
				if got != want {
					t.Errorf("legacy scope %q exposes %q = %t, want %t", test.scope, name, got, want)
				}
			}
		})
	}
}

func TestToolAccessScopesSpecialCases(t *testing.T) {
	tests := []struct {
		name   string
		scopes []oauthentity.Scope
		tool   string
		want   bool
	}{
		{name: "commands read allows list", scopes: []oauthentity.Scope{"commands:read"}, tool: "list_commands", want: true},
		{name: "commands read rejects edit", scopes: []oauthentity.Scope{"commands:read"}, tool: "create_command", want: false},
		{name: "commands edit allows read", scopes: []oauthentity.Scope{"commands:edit"}, tool: "list_commands", want: true},
		{name: "commands edit allows edit", scopes: []oauthentity.Scope{"commands:edit"}, tool: "create_command", want: true},
		{name: "cross group isolation", scopes: []oauthentity.Scope{"commands:edit"}, tool: "list_timers", want: false},
		{name: "overlay read rejects list", scopes: []oauthentity.Scope{"overlays:read"}, tool: "list_overlays", want: false},
		{name: "overlay read rejects get", scopes: []oauthentity.Scope{"overlays:read"}, tool: "get_overlay", want: false},
		{name: "secrets read allows metadata", scopes: []oauthentity.Scope{"secrets:read"}, tool: "list_secrets", want: true},
		{name: "secrets read rejects decrypted value", scopes: []oauthentity.Scope{"secrets:read"}, tool: "get_secret", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			scopes, ok := toolAccessScopesFromOAuthScopes(test.scopes)
			if !ok {
				t.Fatalf("scopes %v did not map", test.scopes)
			}

			// When
			got := scopes.allowsTool(test.tool)

			// Then
			if got != test.want {
				t.Fatalf("allowsTool(%q) = %t, want %t", test.tool, got, test.want)
			}
		})
	}
}

func TestToolAccessScopesFromOAuthScopesRejectsMalformedGrants(t *testing.T) {
	for _, scopes := range [][]oauthentity.Scope{
		nil,
		{"unknown:read"},
		{"commands:unknown"},
		{"commands"},
	} {
		_, ok := toolAccessScopesFromOAuthScopes(scopes)
		if ok {
			t.Fatalf("scopes %v unexpectedly mapped to tool access", scopes)
		}
	}
}

func TestAddToolPanicsWhenPermissionIsMissing(t *testing.T) {
	// Given
	server := modelsdk.NewServer(&modelsdk.Implementation{Name: "test", Version: "1.0.0"}, nil)
	registrar := newToolRegistrar(server, fullToolAccessScopes())

	defer func() {
		// Then
		if recover() == nil {
			t.Fatal("unmapped tool registration did not panic")
		}
	}()

	// When
	addTool(registrar, &modelsdk.Tool{Name: "unmapped_tool"}, func(context.Context, *modelsdk.CallToolRequest, struct{}) (*modelsdk.CallToolResult, any, error) {
		return nil, nil, nil
	})
}

func expectedPermissionMap() map[string]toolPermission {
	permissions := make(map[string]toolPermission, registeredToolCount)
	for _, group := range expectedToolPermissions {
		for name := range strings.FieldsSeq(group.read) {
			permissions[name] = toolPermission{Group: group.group, Action: oauthentity.ScopeActionRead}
		}
		for name := range strings.FieldsSeq(group.edit) {
			permissions[name] = toolPermission{Group: group.group, Action: oauthentity.ScopeActionEdit}
		}
	}
	return permissions
}
