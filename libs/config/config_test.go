package cfg

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewWithEnvPath_ValidatesVKVideoConfigurationOnlyWhenEnabled(t *testing.T) {
	setRequiredConfigEnv(t)
	envPath := filepath.Join(t.TempDir(), "missing.env")

	if _, err := NewWithEnvPath(envPath); err != nil {
		t.Fatalf("configuration without VK Video credentials returned an error: %v", err)
	}

	t.Setenv("VK_VIDEO_CLIENT_ID", "vk-client-id")
	t.Setenv("VK_VIDEO_CLIENT_SECRET", "vk-client-secret")
	if _, err := NewWithEnvPath(envPath); err == nil {
		t.Fatal("VK Video credentials without webhook secret must fail")
	}
}

func TestNewWithEnvPath_LoadsVKVideoConfiguration(t *testing.T) {
	setRequiredConfigEnv(t)
	t.Setenv("VK_VIDEO_CLIENT_ID", "vk-client-id")
	t.Setenv("VK_VIDEO_CLIENT_SECRET", "vk-client-secret")
	t.Setenv("VK_VIDEO_WEBHOOK_SECRET", "vk-webhook-secret")
	t.Setenv("VK_VIDEO_API_BASE_URL", "https://api.example.test")
	t.Setenv("VK_VIDEO_AUTH_BASE_URL", "https://auth.example.test")
	t.Setenv("VK_VIDEO_DEVAPI_BASE_URL", "https://devapi.example.test")

	config, err := NewWithEnvPath(filepath.Join(t.TempDir(), "missing.env"))
	if err != nil {
		t.Fatalf("load VK Video configuration: %v", err)
	}

	want := map[string]any{
		"VKVideoClientID":      "vk-client-id",
		"VKVideoClientSecret":  "vk-client-secret",
		"VKVideoWebhookSecret": "vk-webhook-secret",
		"VKVideoAPIBaseURL":    "https://api.example.test",
		"VKVideoAuthBaseURL":   "https://auth.example.test",
		"VKVideoDevAPIBaseURL": "https://devapi.example.test",
	}

	value := reflect.ValueOf(config).Elem()
	for fieldName, expected := range want {
		field := value.FieldByName(fieldName)
		if !field.IsValid() {
			t.Errorf("Config is missing %s", fieldName)
			continue
		}

		if got := field.Interface(); got != expected {
			t.Errorf("Config.%s = %#v, want %#v", fieldName, got, expected)
		}
	}

	if !config.IsVkVideoEnabled() {
		t.Error("IsVkVideoEnabled() = false, want true with client ID and secret set")
	}
}

func setRequiredConfigEnv(t *testing.T) {
	t.Helper()
	t.Setenv("TWITCH_CLIENTID", "twitch-client-id")
	t.Setenv("TWITCH_CLIENTSECRET", "twitch-client-secret")
	t.Setenv("DATABASE_URL", "postgres://twir:twir@localhost:5432/twir")
}
