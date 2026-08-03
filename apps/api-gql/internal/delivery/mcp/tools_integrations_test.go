package mcp

import (
	"errors"
	"testing"
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

func TestIntegrationStatusRetainsDetailsWhenScopeAllowsWrites(t *testing.T) {
	// Given
	requestScope := scope{AccessScopes: fullToolAccessScopes()}
	details := integrationStatusDetails{AccessToken: "access-token", RefreshToken: "refresh-token"}

	// When
	result := integrationStatus(requestScope, details, nil)

	// Then
	if result.Data != details {
		t.Fatalf("details = %#v, want %#v", result.Data, details)
	}
}
