package mcp_oauth

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestRefreshTokenReuseError_matches_the_typed_reuse_sentinel(t *testing.T) {
	// Given
	wantFamilyID := uuid.New()
	err := &RefreshTokenReuseError{FamilyID: wantFamilyID}

	// When
	var got *RefreshTokenReuseError
	matched := errors.Is(err, ErrRefreshTokenReuse)
	extracted := errors.As(err, &got)

	// Then
	if !matched {
		t.Fatal("reuse error must match ErrRefreshTokenReuse")
	}
	if !extracted || got.FamilyID != wantFamilyID {
		t.Fatalf("reuse error = %#v, want family ID %s", got, wantFamilyID)
	}
}
