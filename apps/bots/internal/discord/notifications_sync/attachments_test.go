package notifications_sync

import (
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
)

func TestMediaSourcesDeduplicatesEmbedImages(t *testing.T) {
	message := discord.Message{
		Attachments: []discord.Attachment{
			{
				ID:          1,
				Filename:    "release.png",
				ContentType: "image/png",
				URL:         "https://cdn.example/release.png",
			},
		},
		Embeds: []discord.Embed{
			{
				Image: &discord.EmbedImage{URL: "https://cdn.example/release.png"},
				Thumbnail: &discord.EmbedThumbnail{
					URL: "https://cdn.example/thumbnail.webp",
				},
			},
		},
	}

	sources := mediaSources(message)
	if len(sources) != 2 {
		t.Fatalf("expected 2 unique media sources, got %d", len(sources))
	}
	if sources[0].Filename != "release.png" {
		t.Fatalf("unexpected attachment filename %q", sources[0].Filename)
	}
	if sources[1].Filename != "thumbnail.webp" {
		t.Fatalf("unexpected embed filename %q", sources[1].Filename)
	}
}

func TestSafeFilename(t *testing.T) {
	if value := safeFilename("../../release notes (final).png"); value != "release_notes_final_.png" {
		t.Fatalf("unexpected safe filename %q", value)
	}
}

func TestDiscordMediaURLAllowlist(t *testing.T) {
	tests := []struct {
		url     string
		allowed bool
	}{
		{url: "https://cdn.discordapp.com/attachments/1/2/image.png", allowed: true},
		{url: "https://media.discordapp.net/attachments/1/2/image.png", allowed: true},
		{url: "http://cdn.discordapp.com/attachments/1/2/image.png", allowed: false},
		{url: "https://cdn.discordapp.com.evil.example/image.png", allowed: false},
		{url: "https://127.0.0.1/internal.png", allowed: false},
	}

	for _, test := range tests {
		if actual := isDiscordMediaURL(test.url); actual != test.allowed {
			t.Errorf("isDiscordMediaURL(%q) = %v, want %v", test.url, actual, test.allowed)
		}
	}
}
