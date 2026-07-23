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
		"type ChannelPlatformOption",
		"channelPlatformBindings",
		"channelPlatformOptions",
		"channelPlatformConnect",
		"channelPlatformDisconnect",
		"channelPlatformSetEnabled",
	} {
		if !strings.Contains(schema.String(), want) {
			t.Errorf("schema does not declare %q", want)
		}
	}

	normalizedSchema := strings.Join(strings.Fields(schema.String()), " ")
	for _, want := range []string{
		"channelPlatformBindings: [ChannelPlatformBinding!]! @isAuthenticated @hasAccessToSelectedDashboard @hasChannelRolesDashboardPermission(permission: VIEW_BOT_SETTINGS)",
		"channelPlatformOptions: [ChannelPlatformOption!]! @isAuthenticated @hasAccessToSelectedDashboard @hasChannelRolesDashboardPermission(permission: VIEW_BOT_SETTINGS)",
		"channelPlatformConnect(platform: Platform!): String! @isAuthenticated @hasAccessToSelectedDashboard @hasChannelRolesDashboardPermission(permission: MANAGE_BOT_SETTINGS)",
		"channelPlatformDisconnect(platform: Platform!): Boolean! @isAuthenticated @hasAccessToSelectedDashboard @hasChannelRolesDashboardPermission(permission: MANAGE_BOT_SETTINGS)",
		"channelPlatformSetEnabled(platform: Platform!, enabled: Boolean!): ChannelPlatformBinding! @isAuthenticated @hasAccessToSelectedDashboard @hasChannelRolesDashboardPermission(permission: MANAGE_BOT_SETTINGS)",
	} {
		if !strings.Contains(normalizedSchema, want) {
			t.Errorf("schema does not declare %q", want)
		}
	}
	if strings.Contains(normalizedSchema, "channelPlatformConnect(platform: Platform!, redirectTo:") {
		t.Error("channelPlatformConnect exposes a client-controlled redirect")
	}
}
