package cfg

import "testing"

func TestConfigVKVideoCallbacksUseSeparateRoutes(t *testing.T) {
	// Given
	config := Config{SiteBaseUrl: "https://twir.example.test"}

	// When
	broadcasterCallbackURL := config.GetVkCallbackUrl()
	botCallbackURL := config.GetVkVideoBotCallbackUrl()

	// Then
	if broadcasterCallbackURL != "https://twir.example.test/login/vk" {
		t.Errorf("broadcaster callback URL = %q", broadcasterCallbackURL)
	}
	if botCallbackURL != "https://twir.example.test/api/auth/vk-video/bot-callback" {
		t.Errorf("bot callback URL = %q", botCallbackURL)
	}
}
