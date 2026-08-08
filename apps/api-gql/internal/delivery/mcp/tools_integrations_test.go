package mcp

import (
	"errors"
	"testing"

	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
)

type integrationStatusDetails struct {
	AccessToken  string
	RefreshToken string
}

func TestIntegrationStatusOmitsDetailsWhenScopeIsReadOnly(t *testing.T) {
	// Given
	requestScope := scope{AccessScopes: toolAccessScopes{toolAccessScopeRead: {}}}
	details := integrationStatusDetails{AccessToken: "access-token", RefreshToken: "refresh-token"}

	// When
	result := integrationStatus(requestScope, details, nil)

	// Then
	if !result.Connected {
		t.Fatal("read status did not preserve connected state")
	}
	if result.Data != nil {
		t.Fatalf("read status exposed details: %#v", result.Data)
	}
}

func TestIntegrationStatusPreservesErrorsWhenScopeIsReadOnly(t *testing.T) {
	// Given
	requestScope := scope{AccessScopes: toolAccessScopes{toolAccessScopeRead: {}}}
	err := errors.New("integration unavailable")

	// When
	result := integrationStatus(requestScope, integrationStatusDetails{}, err)

	// Then
	if result.Error != err.Error() {
		t.Fatalf("error = %q, want %q", result.Error, err.Error())
	}
	if result.Data != nil {
		t.Fatalf("read error status exposed details: %#v", result.Data)
	}
}

func TestIntegrationStatusOmitsDetailsWithoutIntegrationEditScope(t *testing.T) {
	for _, test := range []struct {
		name   string
		scopes toolAccessScopes
	}{
		{name: "integration read", scopes: toolAccessScopes{entity.Scope("integrations:read"): {}}},
		{name: "unrelated edit", scopes: toolAccessScopes{entity.Scope("secrets:edit"): {}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Given
			requestScope := scope{AccessScopes: test.scopes}
			details := integrationStatusDetails{AccessToken: "access-token", RefreshToken: "refresh-token"}

			// When
			result := integrationStatus(requestScope, details, nil)

			// Then
			if result.Data != nil {
				t.Fatalf("scope %s exposed details: %#v", test.name, result.Data)
			}
		})
	}
}

func TestIntegrationStatusRetainsDetailsWithIntegrationEditScope(t *testing.T) {
	// Given
	requestScope := scope{AccessScopes: toolAccessScopes{entity.Scope("integrations:edit"): {}}}
	details := integrationStatusDetails{AccessToken: "access-token", RefreshToken: "refresh-token"}

	// When
	result := integrationStatus(requestScope, details, nil)

	// Then
	if result.Data != details {
		t.Fatalf("details = %#v, want %#v", result.Data, details)
	}
}

func TestLegacyIntegrationResultPreservesDisabledState(t *testing.T) {
	readResult := legacyIntegrationResult(scope{AccessScopes: toolAccessScopes{toolAccessScopeRead: {}}}, false)
	if readResult.Connected || readResult.Data != nil {
		t.Fatalf("read result = %#v, want disconnected without data", readResult)
	}

	writeResult := legacyIntegrationResult(scope{AccessScopes: toolAccessScopes{entity.Scope("integrations:edit"): {}}}, false)
	if writeResult.Connected {
		t.Fatalf("write result = %#v, want disconnected", writeResult)
	}
	data, ok := writeResult.Data.(map[string]any)
	if !ok || data["enabled"] != false {
		t.Fatalf("write data = %#v, want enabled=false", writeResult.Data)
	}
}
