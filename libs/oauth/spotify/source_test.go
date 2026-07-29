package spotify

import (
	"errors"
	"testing"

	"github.com/twirapp/twir/libs/oauth"
)

func TestNewTokenSource_rejects_missing_dependencies(t *testing.T) {
	// Given
	options := SourceOptions{}

	// When
	_, err := NewTokenSource(options, nil, nil)

	// Then
	if !errors.Is(err, oauth.ErrInvalidOption) {
		t.Fatalf("error = %v", err)
	}
}
