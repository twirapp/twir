package channel_integrations

import (
	"context"
	"errors"
	"testing"

	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
)

func TestSource_returns_typed_error_for_unsupported_service(t *testing.T) {
	// Given
	source := New(nil, nil)

	// When
	_, err := source.Token(context.Background(), integrationsmodel.ServiceLastfm, "channel-1")

	// Then
	if !errors.Is(err, ErrUnsupportedService) {
		t.Fatalf("error = %v", err)
	}
}
