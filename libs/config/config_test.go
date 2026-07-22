package cfg

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewWithEnvPath_ValidatesVKVideoConfigurationOnlyWhenEnabled(t *testing.T) {
	setRequiredConfigEnv(t)
	envPath := filepath.Join(t.TempDir(), "missing.env")

	t.Setenv("VK_VIDEO_ENABLED", "false")
	if _, err := NewWithEnvPath(envPath); err != nil {
		t.Fatalf("disabled VK Video configuration returned an error: %v", err)
	}

	t.Setenv("VK_VIDEO_ENABLED", "true")
	if _, err := NewWithEnvPath(envPath); err == nil {
		t.Fatal("enabled VK Video configuration without credentials must fail")
	}
}

func TestNewWithEnvPath_LoadsVKVideoConfiguration(t *testing.T) {
	setRequiredConfigEnv(t)
	t.Setenv("VK_VIDEO_ENABLED", "true")
	t.Setenv("VK_VIDEO_CLIENT_ID", "vk-client-id")
	t.Setenv("VK_VIDEO_CLIENT_SECRET", "vk-client-secret")
	t.Setenv("VK_VIDEO_SERVICE_TOKEN", "vk-service-token")
	t.Setenv("VK_VIDEO_CALLBACK_URL", "https://twir.example.test/auth/vk/callback")
	t.Setenv("VK_VIDEO_WEBHOOK_SECRET", "vk-webhook-secret")
	t.Setenv("VK_VIDEO_API_BASE_URL", "https://id.example.test")

	config, err := NewWithEnvPath(filepath.Join(t.TempDir(), "missing.env"))
	if err != nil {
		t.Fatalf("load VK Video configuration: %v", err)
	}

	want := map[string]any{
		"VKVideoEnabled":       true,
		"VKVideoClientID":      "vk-client-id",
		"VKVideoClientSecret":  "vk-client-secret",
		"VKVideoServiceToken":  "vk-service-token",
		"VKVideoCallbackURL":   "https://twir.example.test/auth/vk/callback",
		"VKVideoWebhookSecret": "vk-webhook-secret",
		"VKVideoAPIBaseURL":    "https://id.example.test",
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
}

func setRequiredConfigEnv(t *testing.T) {
	t.Helper()
	t.Setenv("TWITCH_CLIENTID", "twitch-client-id")
	t.Setenv("TWITCH_CLIENTSECRET", "twitch-client-secret")
	t.Setenv("DATABASE_URL", "postgres://twir:twir@localhost:5432/twir")
}
