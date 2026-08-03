package mcp_oauth

import (
	"errors"
	"reflect"
	"testing"
)

var expectedScopeGroups = []ScopeGroup{
	"commands", "timers", "files", "games", "song_requests", "moderation",
	"overlays", "integrations", "events", "rewards", "giveaways", "greetings",
	"notifications", "alerts", "secrets", "storage", "pastes", "short_urls",
	"dashboard", "variables", "quotes", "keywords",
}

var expectedScopeStrings = []string{
	"commands:read", "commands:edit", "timers:read", "timers:edit",
	"files:read", "files:edit", "games:read", "games:edit",
	"song_requests:read", "song_requests:edit", "moderation:read", "moderation:edit",
	"overlays:read", "overlays:edit", "integrations:read", "integrations:edit",
	"events:read", "events:edit", "rewards:read", "rewards:edit",
	"giveaways:read", "giveaways:edit", "greetings:read", "greetings:edit",
	"notifications:read", "notifications:edit", "alerts:read", "alerts:edit",
	"secrets:read", "secrets:edit", "storage:read", "storage:edit",
	"pastes:read", "pastes:edit", "short_urls:read", "short_urls:edit",
	"dashboard:read", "dashboard:edit", "variables:read", "variables:edit",
	"quotes:read", "quotes:edit", "keywords:read", "keywords:edit",
}

func TestAllScopeGroups_returns_canonical_catalog_and_clone(t *testing.T) {
	// Given
	wantGroups := expectedScopeGroups

	// When
	gotGroups := AllScopeGroups()

	// Then
	if len(gotGroups) != len(wantGroups) {
		t.Fatalf("group count: got %d, want %d", len(gotGroups), len(wantGroups))
	}
	for index, group := range gotGroups {
		if group.Group != wantGroups[index] {
			t.Errorf("group %d: got %q, want %q", index, group.Group, wantGroups[index])
		}
		if group.Name == "" || group.Description == "" {
			t.Errorf("group %q has incomplete catalog metadata", group.Group)
		}
	}
	if gotGroups[0].Name != "Commands" {
		t.Errorf("commands name: got %q, want %q", gotGroups[0].Name, "Commands")
	}
	if gotGroups[0].Description != "View and manage custom commands, groups, and role cooldowns" {
		t.Errorf("commands description: got %q", gotGroups[0].Description)
	}

	gotGroups[0].Group = "changed"
	clonedGroups := AllScopeGroups()
	if clonedGroups[0].Group != wantGroups[0] {
		t.Errorf("catalog shares mutable state: got %q, want %q", clonedGroups[0].Group, wantGroups[0])
	}
}

func TestAllScopes_returns_canonical_order_and_clone(t *testing.T) {
	// Given
	want := scopesFromStrings(expectedScopeStrings)

	// When
	got := AllScopes()

	// Then
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("scopes: got %v, want %v", got, want)
	}
	got[0] = Scope("changed")
	if fresh := AllScopes(); fresh[0] != want[0] {
		t.Errorf("scopes share mutable state: got %q, want %q", fresh[0], want[0])
	}
}

func TestParseScopes_normalizes_legacy_granular_and_duplicate_tokens(t *testing.T) {
	// Given
	all := scopesFromStrings(expectedScopeStrings)
	readOnly := make([]Scope, 0, len(expectedScopeGroups))
	for _, group := range expectedScopeGroups {
		readOnly = append(readOnly, Scope(string(group)+":read"))
	}
	readAndTimerEdit := append([]Scope{"commands:read", "timers:read", "timers:edit"}, readOnly[2:]...)

	tests := []struct {
		name string
		raw  string
		want []Scope
	}{
		{name: "granular read", raw: "commands:read", want: scopesFromStrings([]string{"commands:read"})},
		{name: "granular edit implies read", raw: "commands:edit", want: scopesFromStrings([]string{"commands:read", "commands:edit"})},
		{name: "legacy read expands all groups", raw: "read", want: readOnly},
		{name: "legacy write expands all groups and reads", raw: "write", want: all},
		{name: "legacy read and write produce full catalog", raw: "read write", want: all},
		{name: "mixed legacy and granular tokens union", raw: "read timers:edit", want: readAndTimerEdit},
		{name: "duplicates are removed", raw: "commands:read commands:read commands:edit", want: scopesFromStrings([]string{"commands:read", "commands:edit"})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			got, err := ParseScopes(tt.raw)

			// Then
			if err != nil {
				t.Fatalf("ParseScopes() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseScopes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseScopes_rejects_invalid_tokens_with_sentinel(t *testing.T) {
	// Given
	inputs := []string{"", " ", "unknown:read", "commands:unknown", "commands", "commands:read:extra", ":read", "commands:", "commands::read", "Commands:read", "commands:Read"}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			// When
			_, err := ParseScopes(input)

			// Then
			if err == nil || !errors.Is(err, ErrInvalidScope) {
				t.Errorf("ParseScopes(%q) error = %v, want wrapped ErrInvalidScope", input, err)
			}
		})
	}
}

func TestNormalizeScopes_rejects_empty_and_invalid_values(t *testing.T) {
	// Given
	tests := []struct {
		name    string
		input   []Scope
		wantErr bool
	}{
		{name: "empty", input: nil, wantErr: true},
		{name: "unknown group", input: []Scope{"unknown:read"}, wantErr: true},
		{name: "unknown action", input: []Scope{"commands:unknown"}, wantErr: true},
		{name: "missing colon", input: []Scope{"commands"}, wantErr: true},
		{name: "extra colon", input: []Scope{"commands:read:extra"}, wantErr: true},
		{name: "empty segment", input: []Scope{"commands:"}, wantErr: true},
		{name: "case variant", input: []Scope{"COMMANDS:read"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			_, err := NormalizeScopes(tt.input)

			// Then
			if tt.wantErr && (err == nil || !errors.Is(err, ErrInvalidScope)) {
				t.Errorf("NormalizeScopes() error = %v, want wrapped ErrInvalidScope", err)
			}
		})
	}
}

func TestNormalizeScopes_returns_canonical_union(t *testing.T) {
	// Given
	input := []Scope{"timers:edit", "commands:read", "timers:edit"}

	// When
	got, err := NormalizeScopes(input)

	// Then
	if err != nil {
		t.Fatalf("NormalizeScopes() error = %v", err)
	}
	want := scopesFromStrings([]string{"commands:read", "timers:read", "timers:edit"})
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NormalizeScopes() = %v, want %v", got, want)
	}
}

func TestScopeSubset_normalizes_both_operands(t *testing.T) {
	// Given
	tests := []struct {
		name      string
		requested []Scope
		allowed   []Scope
		want      bool
	}{
		{name: "edit allows read", requested: []Scope{"commands:read"}, allowed: []Scope{"commands:edit"}, want: true},
		{name: "read does not allow edit", requested: []Scope{"commands:edit"}, allowed: []Scope{"commands:read"}, want: false},
		{name: "legacy read requested", requested: []Scope{ScopeRead}, allowed: []Scope{"commands:edit"}, want: false},
		{name: "legacy write allowed", requested: []Scope{"commands:edit"}, allowed: []Scope{ScopeWrite}, want: true},
		{name: "invalid requested scope", requested: []Scope{"commands:bad"}, allowed: []Scope{ScopeWrite}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			got := ScopeSubset(tt.requested, tt.allowed)

			// Then
			if got != tt.want {
				t.Errorf("ScopeSubset() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestHasScope_uses_normalized_scope_set(t *testing.T) {
	// Given
	scopes := []Scope{"commands:edit", "timers:read"}

	// When
	commandsRead := HasScope(scopes, ScopeGroup("commands"), ScopeAction("read"))
	commandsEdit := HasScope(scopes, ScopeGroup("commands"), ScopeAction("edit"))
	filesRead := HasScope([]Scope{ScopeRead}, ScopeGroup("files"), ScopeAction("read"))
	filesEdit := HasScope([]Scope{ScopeRead}, ScopeGroup("files"), ScopeAction("edit"))

	// Then
	if !commandsRead || !commandsEdit || !filesRead || filesEdit {
		t.Errorf("HasScope() = commands read %t, commands edit %t, files read %t, files edit %t", commandsRead, commandsEdit, filesRead, filesEdit)
	}
}

func TestScopeStrings_returns_a_fresh_string_slice(t *testing.T) {
	// Given
	scopes := []Scope{"commands:read", "commands:edit"}

	// When
	got := ScopeStrings(scopes)
	got[0] = "changed"

	// Then
	if scopes[0] != Scope("commands:read") {
		t.Errorf("ScopeStrings() mutated input: got %q", scopes[0])
	}
	if fresh := ScopeStrings(scopes); fresh[0] != "commands:read" {
		t.Errorf("ScopeStrings() shares mutable state: got %q", fresh[0])
	}
}

func scopesFromStrings(values []string) []Scope {
	result := make([]Scope, len(values))
	for index, value := range values {
		result[index] = Scope(value)
	}
	return result
}
