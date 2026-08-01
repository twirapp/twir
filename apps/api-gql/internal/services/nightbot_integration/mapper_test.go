package nightbot_integration

import (
	"reflect"
	"testing"

	"github.com/twirapp/twir/apps/api-gql/internal/services/importer"
)

func TestNormalizeCommandsConvertsSupportedNightbotCommands(t *testing.T) {
	t.Parallel()

	ownerAlias := "!owner"
	commands := nightbotCustomCommandsResponse{
		Commands: []struct {
			Alias     *string `json:"alias,omitempty"`
			Name      string  `json:"name"`
			Message   string  `json:"message"`
			UserLevel string  `json:"userLevel"`
			CoolDown  int     `json:"coolDown"`
			Count     int     `json:"count"`
		}{
			{Name: "!Owner", Message: "owner response", UserLevel: "owner", CoolDown: 10},
			{Name: "!OwnerAlias", Alias: &ownerAlias, Message: "owner alias", UserLevel: "owner"},
			{Name: "!Moderator", Message: "moderator response", UserLevel: "moderator"},
			{Name: "!Subscriber", Message: "subscriber response", UserLevel: "subscriber"},
			{Name: "!VIP", Message: "vip response", UserLevel: "twitch_vip"},
			{Name: "!Everyone", Message: "everyone response", UserLevel: "everyone"},
		},
	}

	got, failures := NormalizeCommands(commands)
	if want := []importer.Command{
		{
			Name: "owner", Response: "owner response", Enabled: true, Visible: true, IsReply: true,
			Aliases: []string{"owneralias"}, Cooldown: 10, Role: importer.RoleBroadcaster,
		},
		{
			Name: "moderator", Response: "moderator response", Enabled: true, Visible: true, IsReply: true,
			Aliases: []string{}, Role: importer.RoleModerator,
		},
		{
			Name: "subscriber", Response: "subscriber response", Enabled: true, Visible: true, IsReply: true,
			Aliases: []string{}, Role: importer.RoleSubscriber,
		},
		{
			Name: "vip", Response: "vip response", Enabled: true, Visible: true, IsReply: true,
			Aliases: []string{}, Role: importer.RoleVip,
		},
		{
			Name: "everyone", Response: "everyone response", Enabled: true, Visible: true, IsReply: true,
			Aliases: []string{}, Role: importer.RoleEveryone,
		},
	}; !reflect.DeepEqual(got, want) {
		t.Fatalf("NormalizeCommands() commands = %#v, want %#v", got, want)
	}
	if want := []importer.Failure{}; !reflect.DeepEqual(failures, want) {
		t.Fatalf("NormalizeCommands() failures = %#v, want %#v", failures, want)
	}
}

func TestNormalizeCommandsReportsUnsupportedNightbotRoles(t *testing.T) {
	t.Parallel()

	commands := nightbotCustomCommandsResponse{
		Commands: []struct {
			Alias     *string `json:"alias,omitempty"`
			Name      string  `json:"name"`
			Message   string  `json:"message"`
			UserLevel string  `json:"userLevel"`
			CoolDown  int     `json:"coolDown"`
			Count     int     `json:"count"`
		}{
			{Name: "!Admin", Message: "admin response", UserLevel: "admin"},
			{Name: "!Regular", Message: "regular response", UserLevel: "regular"},
		},
	}

	got, failures := NormalizeCommands(commands)
	if want := []importer.Command{}; !reflect.DeepEqual(got, want) {
		t.Fatalf("NormalizeCommands() commands = %#v, want %#v", got, want)
	}
	if want := []importer.Failure{
		{Name: "!Admin", Reason: importer.FailureUnsupportedRole},
		{Name: "!Regular", Reason: importer.FailureUnsupportedRole},
	}; !reflect.DeepEqual(failures, want) {
		t.Fatalf("NormalizeCommands() failures = %#v, want %#v", failures, want)
	}
}

func TestNormalizeTimersConvertsCronSchedulesToMinutes(t *testing.T) {
	t.Parallel()

	timers := nightbotTimersResponse{
		Timers: []struct {
			ID       string `json:"_id"`
			Name     string `json:"name"`
			Message  string `json:"message"`
			Interval string `json:"interval"`
			Lines    int    `json:"lines"`
			Enabled  bool   `json:"enabled"`
		}{
			{Name: "five", Message: "every five minutes", Interval: "*/5 * * * *", Lines: 2, Enabled: true},
			{Name: "hourly", Message: "once an hour", Interval: "15 * * * *", Lines: 3, Enabled: false},
		},
	}

	got, failures := NormalizeTimers(timers)
	if want := []importer.Timer{
		{
			Name: "five", Message: "every five minutes", Enabled: true, OnlineEnabled: true,
			TimeInterval: 5, MessageInterval: 2,
		},
		{
			Name: "hourly", Message: "once an hour", Enabled: false, OnlineEnabled: true,
			TimeInterval: 60, MessageInterval: 3,
		},
	}; !reflect.DeepEqual(got, want) {
		t.Fatalf("NormalizeTimers() timers = %#v, want %#v", got, want)
	}
	if want := []importer.Failure{}; !reflect.DeepEqual(failures, want) {
		t.Fatalf("NormalizeTimers() failures = %#v, want %#v", failures, want)
	}
}

func TestNormalizeTimersReportsInvalidCronSchedules(t *testing.T) {
	t.Parallel()

	timers := nightbotTimersResponse{
		Timers: []struct {
			ID       string `json:"_id"`
			Name     string `json:"name"`
			Message  string `json:"message"`
			Interval string `json:"interval"`
			Lines    int    `json:"lines"`
			Enabled  bool   `json:"enabled"`
		}{
			{Name: "daily", Message: "once a day", Interval: "0 12 * * *", Lines: 1, Enabled: true},
		},
	}

	got, failures := NormalizeTimers(timers)
	if want := []importer.Timer{}; !reflect.DeepEqual(got, want) {
		t.Fatalf("NormalizeTimers() timers = %#v, want %#v", got, want)
	}
	if want := []importer.Failure{{Name: "daily", Reason: importer.FailureIncompatibleInterval}}; !reflect.DeepEqual(failures, want) {
		t.Fatalf("NormalizeTimers() failures = %#v, want %#v", failures, want)
	}
}
