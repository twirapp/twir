package vkvideoprobe

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRedactFrame_redactsKnownCredentialsSensitiveFieldsURLsAndJWTs(t *testing.T) {
	// Given
	raw := []byte(`{"authorization":"Bearer synthetic-access-token","token":"synthetic-connection-token","nested":{"refresh_token":"synthetic-subscription-token"},"url":"wss://example.test/ws?access_token=synthetic-access-token","jwt":"eyJheader.payload.signature"}`)

	// When
	got := string(RedactFrame(raw, []string{"synthetic-access-token", "synthetic-connection-token", "synthetic-subscription-token"}))

	// Then
	for _, secret := range []string{"synthetic-access-token", "synthetic-connection-token", "synthetic-subscription-token", "eyJheader.payload.signature"} {
		if strings.Contains(got, secret) {
			t.Fatalf("redacted frame still contains %q: %s", secret, got)
		}
	}
	if !strings.Contains(got, "[REDACTED]") {
		t.Fatalf("redacted frame = %s", got)
	}
}

func TestCreateOutput_refusesExistingFile(t *testing.T) {
	// Given
	path := filepath.Join(t.TempDir(), "capture.jsonl")
	if err := os.WriteFile(path, []byte("existing"), 0600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	// When
	file, err := CreateOutput(path)

	// Then
	if file != nil {
		_ = file.Close()
		t.Fatal("output file should not have been opened")
	}
	if !errors.Is(err, ErrOutputExists) {
		t.Fatalf("create output error = %v, want ErrOutputExists", err)
	}
}
