package streamelements

import (
	"reflect"
	"testing"

	"github.com/twirapp/twir/apps/api-gql/internal/services/importer"
	streamelementsintegration "github.com/twirapp/twir/libs/integrations/streamelements"
)

func TestNormalizeCommandsConvertsSupportedStreamElementsCommands(t *testing.T) {
	t.Parallel()

	input := []streamelementsintegration.Command{
		{
			Name: "everyone", Response: "everyone response", Enabled: true,
			AccessLevel: 100, Aliases: []string{"hi", "hello"},
			Cooldown:      streamelementsintegration.CommandCooldown{Global: 15},
			EnabledOnline: true, Type: "reply",
		},
		{Name: "subscriber", Response: "subscriber response", AccessLevel: 250, EnabledOffline: true, Type: "say"},
		{Name: "vip", Response: "vip response", AccessLevel: 400, EnabledOnline: true, EnabledOffline: true, Type: "action"},
		{Name: "moderator", Response: "moderator response", AccessLevel: 500, Hidden: true, Type: "say"},
		{Name: "owner", Response: "owner response", AccessLevel: 1000, Type: "reply"},
		{Name: "owner-plus", Response: "owner plus response", AccessLevel: 1500, Type: "reply"},
	}

	got, failures := NormalizeCommands(input)
	want := []importer.Command{
		{
			Name: "everyone", Response: "everyone response", Enabled: true, Visible: true, IsReply: true,
			Aliases: []string{"hi", "hello"}, Cooldown: 15, Role: importer.RoleEveryone, OnlineOnly: true,
		},
		{
			Name: "subscriber", Response: "subscriber response", Visible: true,
			Aliases: []string{}, Role: importer.RoleSubscriber, OfflineOnly: true,
		},
		{
			Name: "vip", Response: "/me vip response", Visible: true,
			Aliases: []string{}, Role: importer.RoleVip,
		},
		{
			Name: "moderator", Response: "moderator response",
			Aliases: []string{}, Role: importer.RoleModerator,
		},
		{
			Name: "owner", Response: "owner response", Visible: true, IsReply: true,
			Aliases: []string{}, Role: importer.RoleBroadcaster,
		},
		{
			Name: "owner-plus", Response: "owner plus response", Visible: true, IsReply: true,
			Aliases: []string{}, Role: importer.RoleBroadcaster,
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("NormalizeCommands() commands = %#v, want %#v", got, want)
	}
	if wantFailures := []importer.Failure{}; !reflect.DeepEqual(failures, wantFailures) {
		t.Fatalf("NormalizeCommands() failures = %#v, want %#v", failures, wantFailures)
	}
}

func TestNormalizeCommandsReportsUnsupportedRolesAndResponseTypes(t *testing.T) {
	t.Parallel()

	input := []streamelementsintegration.Command{
		{Name: "regular", Response: "regular response", AccessLevel: 300, Type: "say"},
		{Name: "unknown-role", Response: "unknown role response", AccessLevel: 999, Type: "say"},
		{Name: "whisper", Response: "whisper response", AccessLevel: 100, Type: "whisper"},
		{Name: "custom", Response: "custom response", AccessLevel: 100, Type: "custom"},
	}

	got, failures := NormalizeCommands(input)
	if want := []importer.Command{}; !reflect.DeepEqual(got, want) {
		t.Fatalf("NormalizeCommands() commands = %#v, want %#v", got, want)
	}
	wantFailures := []importer.Failure{
		{Name: "regular", Reason: importer.FailureUnsupportedRole},
		{Name: "unknown-role", Reason: importer.FailureUnsupportedRole},
		{Name: "whisper", Reason: importer.FailureUnsupportedResponse},
		{Name: "custom", Reason: importer.FailureUnsupportedResponse},
	}
	if !reflect.DeepEqual(failures, wantFailures) {
		t.Fatalf("NormalizeCommands() failures = %#v, want %#v", failures, wantFailures)
	}
}

func TestNormalizeCommandsPreservesActionsAndDisablesCommandsWithNoActiveMode(t *testing.T) {
	t.Parallel()

	input := []streamelementsintegration.Command{
		{
			Name: "action", Response: "waves", Enabled: true, AccessLevel: 100,
			EnabledOnline: true, EnabledOffline: true, Type: "action",
		},
		{
			Name: "no-mode", Response: "hidden response", Enabled: true,
			AccessLevel: 100, Type: "say",
		},
	}

	got, failures := NormalizeCommands(input)
	want := []importer.Command{
		{
			Name: "action", Response: "/me waves", Enabled: true, Visible: true,
			Aliases: []string{}, Role: importer.RoleEveryone,
		},
		{
			Name: "no-mode", Response: "hidden response", Enabled: false, Visible: true,
			Aliases: []string{}, Role: importer.RoleEveryone,
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("NormalizeCommands() commands = %#v, want %#v", got, want)
	}
	if len(failures) != 0 {
		t.Fatalf("NormalizeCommands() failures = %#v, want none", failures)
	}
}

func TestNormalizeTimersConvertsRepresentableModes(t *testing.T) {
	t.Parallel()

	oneMode := streamelementsintegration.Timer{Name: "online", Message: "online response", Enabled: true, ChatLines: 3}
	oneMode.Online.Enabled = true
	oneMode.Online.Interval = 5

	equalModes := streamelementsintegration.Timer{Name: "both", Message: "both response", Enabled: true, ChatLines: 7}
	equalModes.Online.Enabled = true
	equalModes.Online.Interval = 10
	equalModes.Offline.Enabled = true
	equalModes.Offline.Interval = 10

	offlineMode := streamelementsintegration.Timer{Name: "offline", Message: "offline response", ChatLines: 2}
	offlineMode.Offline.Enabled = true
	offlineMode.Offline.Interval = 15

	got, failures := NormalizeTimers([]streamelementsintegration.Timer{oneMode, equalModes, offlineMode})
	want := []importer.Timer{
		{
			Name: "online", Message: "online response", Enabled: true, OnlineEnabled: true,
			TimeInterval: 5, MessageInterval: 3,
		},
		{
			Name: "both", Message: "both response", Enabled: true, OnlineEnabled: true, OfflineEnabled: true,
			TimeInterval: 10, MessageInterval: 7,
		},
		{
			Name: "offline", Message: "offline response", OfflineEnabled: true,
			TimeInterval: 15, MessageInterval: 2,
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("NormalizeTimers() timers = %#v, want %#v", got, want)
	}
	if wantFailures := []importer.Failure{}; !reflect.DeepEqual(failures, wantFailures) {
		t.Fatalf("NormalizeTimers() failures = %#v, want %#v", failures, wantFailures)
	}
}

func TestNormalizeTimersReportsIncompatibleIntervals(t *testing.T) {
	t.Parallel()

	timer := streamelementsintegration.Timer{Name: "different", Message: "response", Enabled: true, ChatLines: 1}
	timer.Online.Enabled = true
	timer.Online.Interval = 5
	timer.Offline.Enabled = true
	timer.Offline.Interval = 10

	got, failures := NormalizeTimers([]streamelementsintegration.Timer{timer})
	if want := []importer.Timer{}; !reflect.DeepEqual(got, want) {
		t.Fatalf("NormalizeTimers() timers = %#v, want %#v", got, want)
	}
	if want := []importer.Failure{{Name: "different", Reason: importer.FailureIncompatibleInterval}}; !reflect.DeepEqual(failures, want) {
		t.Fatalf("NormalizeTimers() failures = %#v, want %#v", failures, want)
	}
}
