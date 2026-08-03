package notifications_sync

import (
	"strings"
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/goccy/go-json"
)

func TestBuildEditorJSPreservesDiscordFormatting(t *testing.T) {
	message := discord.Message{
		ID:      42,
		Content: "# Release\n**Bold** and ~~old~~ with \x60code\x60\n- first\n- second\n> quoted",
	}

	value, ok, err := buildEditorJS(message, []renderedMedia{
		{URL: "https://cdn.example/image.png", Filename: "image.png", ContentType: "image/png", IsImage: true},
	})
	if err != nil {
		t.Fatalf("build editor.js: %v", err)
	}
	if !ok {
		t.Fatal("expected document")
	}

	var document editorDocument
	if err := json.Unmarshal([]byte(value), &document); err != nil {
		t.Fatalf("decode editor.js: %v", err)
	}
	if len(document.Blocks) != 5 {
		t.Fatalf("expected 5 blocks, got %d", len(document.Blocks))
	}
	paragraph, ok := document.Blocks[1].Data.(map[string]any)
	if !ok {
		t.Fatalf("expected paragraph data, got %#v", document.Blocks[1].Data)
	}
	paragraphText, ok := paragraph["text"].(string)
	if !ok {
		t.Fatalf("expected paragraph text, got %#v", paragraph)
	}
	if !strings.Contains(paragraphText, "<strong>Bold</strong>") {
		t.Fatalf("expected bold markup in %s", paragraphText)
	}
	if !strings.Contains(paragraphText, "<s>old</s>") {
		t.Fatalf("expected strikethrough markup in %s", paragraphText)
	}
	if document.Blocks[4].Type != "image" {
		t.Fatalf("expected image block, got %s", document.Blocks[4].Type)
	}
}

func TestBuildEditorJSEscapesHTML(t *testing.T) {
	message := discord.Message{ID: 42, Content: "<script>alert(1)</script>"}

	value, ok, err := buildEditorJS(message, nil)
	if err != nil {
		t.Fatalf("build editor.js: %v", err)
	}
	if !ok {
		t.Fatal("expected document")
	}
	var document editorDocument
	if err := json.Unmarshal([]byte(value), &document); err != nil {
		t.Fatalf("decode editor.js: %v", err)
	}
	data, ok := document.Blocks[0].Data.(map[string]any)
	if !ok {
		t.Fatalf("expected paragraph data, got %#v", document.Blocks[0].Data)
	}
	text, ok := data["text"].(string)
	if !ok {
		t.Fatalf("expected paragraph text, got %#v", data)
	}
	if strings.Contains(text, "<script>") {
		t.Fatalf("unsafe html was not escaped: %s", text)
	}
	if !strings.Contains(text, "&lt;script&gt;") {
		t.Fatalf("escaped html missing: %s", text)
	}
}
