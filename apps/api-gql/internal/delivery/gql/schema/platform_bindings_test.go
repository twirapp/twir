package schema_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestPlatformBindingSchemaDeclaresBindingSurface(t *testing.T) {
	t.Parallel()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("get test file path")
	}

	var schema strings.Builder
	for _, name := range []string{"platform.graphql", "channels.graphql"} {
		contents, err := os.ReadFile(filepath.Join(filepath.Dir(file), name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		schema.Write(contents)
	}

	for _, want := range []string{
		"VK_VIDEO_LIVE",
		"type PlatformCapability",
		"type ChannelPlatformBinding",
		"channelPlatformBindings",
		"channelPlatformConnect",
		"channelPlatformDisconnect",
		"channelPlatformSetEnabled",
	} {
		if !strings.Contains(schema.String(), want) {
			t.Errorf("schema does not declare %q", want)
		}
	}
}
