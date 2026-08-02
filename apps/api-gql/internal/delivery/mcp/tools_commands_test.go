package mcp

import (
	"testing"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/google/uuid"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	"github.com/twirapp/twir/libs/entities/platform"
)

func TestCreateCommandInputSchemaExposesDashboardFields(t *testing.T) {
	schema, err := jsonschema.For[createCommandInput](nil)
	if err != nil {
		t.Fatalf("create input schema: %v", err)
	}

	wantFields := []string{
		"name", "responses", "aliases", "description", "cooldown", "cooldownType",
		"enabled", "visible", "isReply", "keepResponsesOrder", "deniedUsersIds",
		"allowedUsersIds", "rolesIds", "onlineOnly", "offlineOnly", "enabledCategories",
		"requiredWatchTime", "requiredMessages", "requiredUsedChannelPoints", "groupId",
		"expiresAt", "expiresType", "roleCooldowns", "platforms",
	}
	for _, field := range wantFields {
		if schema.Properties[field] == nil {
			t.Errorf("missing field %q from create_command schema", field)
		}
	}
}

func TestCreateCommandServiceInputMapsEverySetting(t *testing.T) {
	channelID := uuid.New()
	groupID := uuid.New().String()
	roleID := uuid.New().String()
	expiresAt := 1_800_000_000_000
	expiresType := "DISABLE"

	got, err := createCommandServiceInput(
		scope{Channel: channelentity.Channel{ID: channelID}, ActorID: "actor"},
		createCommandInput{
			Name:                      "jsonparty",
			Responses:                 []commandResponseInput{{Text: "$(customvar|jsonparty)", Order: 2, TwitchCategoriesIDS: []string{"509658"}, OnlineOnly: true, Platforms: []platform.Platform{platform.PlatformTwitch}}},
			Aliases:                   []string{"json"},
			Description:               "JSON party",
			Cooldown:                  10,
			CooldownType:              "GLOBAL",
			Enabled:                   true,
			Visible:                   true,
			IsReply:                   true,
			KeepResponsesOrder:        true,
			DeniedUsersIDS:            []string{"denied"},
			AllowedUsersIDS:           []string{"allowed"},
			RolesIDS:                  []string{roleID},
			OnlineOnly:                true,
			OfflineOnly:               false,
			RequiredWatchTime:         60,
			RequiredMessages:          20,
			RequiredUsedChannelPoints: 100,
			GroupID:                   &groupID,
			ExpiresAt:                 &expiresAt,
			ExpiresType:               &expiresType,
			RoleCooldowns:             []commandRoleCooldownInput{{RoleID: roleID, Cooldown: 5}},
			Platforms:                 []platform.Platform{platform.PlatformTwitch, platform.PlatformKick},
		},
	)
	if err != nil {
		t.Fatalf("map create input: %v", err)
	}

	if got.ChannelID != channelID.String() || got.ActorID != "actor" || got.GroupID.String() != groupID {
		t.Fatalf("scope fields were not mapped: %#v", got)
	}
	if got.EnabledCategories == nil {
		t.Fatal("enabled categories must be a non-nil empty slice")
	}
	if !got.IsReply || !got.KeepResponsesOrder || !got.OnlineOnly || got.OfflineOnly {
		t.Fatalf("command flags were not mapped: %#v", got)
	}
	if got.RequiredWatchTime != 60 || got.RequiredMessages != 20 || got.RequiredUsedChannelPoints != 100 {
		t.Fatalf("access thresholds were not mapped: %#v", got)
	}
	if len(got.Responses) != 1 || got.Responses[0].Order != 2 || len(got.Responses[0].TwitchCategoryIDs) != 1 || len(got.Responses[0].Platforms) != 1 {
		t.Fatalf("response settings were not mapped: %#v", got.Responses)
	}
	if len(got.RoleCooldowns) != 1 || len(got.Platforms) != 2 || got.ExpiresAt == nil || got.ExpiresType == nil {
		t.Fatalf("advanced settings were not mapped: %#v", got)
	}
}

func TestUpdateCommandServiceInputPreservesOmittedCollections(t *testing.T) {
	got, err := updateCommandServiceInput(scope{Channel: channelentity.Channel{ID: uuid.New()}}, updateCommandInput{})
	if err != nil {
		t.Fatalf("map update input: %v", err)
	}
	if got.Responses != nil || got.RoleCooldowns != nil || got.Platforms != nil || got.EnabledCategories != nil {
		t.Fatalf("omitted collections must remain nil: %#v", got)
	}
}

func TestUpdateCommandServiceInputConvertsExpiration(t *testing.T) {
	expiresAt := 1_800_000_000_000
	expiresType := "DELETE"
	got, err := updateCommandServiceInput(scope{Channel: channelentity.Channel{ID: uuid.New()}}, updateCommandInput{ExpiresAt: &expiresAt, ExpiresType: &expiresType})
	if err != nil {
		t.Fatalf("map update input: %v", err)
	}
	want := time.UnixMilli(int64(expiresAt))
	if got.ExpiresAt == nil || !got.ExpiresAt.Equal(want) {
		t.Fatalf("expiration = %v, want %v", got.ExpiresAt, want)
	}
}

func TestCreateCommandServiceInputRejectsInvalidAdvancedSettings(t *testing.T) {
	badType := "ARCHIVE"
	badID := "not-a-uuid"
	tests := []createCommandInput{
		{ExpiresType: &badType},
		{GroupID: &badID},
		{RolesIDS: []string{badID}},
		{Platforms: []platform.Platform{"discord"}},
	}
	for _, input := range tests {
		if _, err := createCommandServiceInput(scope{Channel: channelentity.Channel{ID: uuid.New()}}, input); err == nil {
			t.Fatalf("expected error for input %#v", input)
		}
	}
}
