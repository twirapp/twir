package dota

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/twirapp/twir/libs/repositories/dota/model"
)

const enabledCommandsSettingsJSON = `{"mmr":true,"wl":true,"lg":true,"gm":true,"np":true,"wp":true}`

func TestCreateInputLeavesCommandSettingsUnspecified(t *testing.T) {
	input := CreateInput{}
	if input.CommandsSettings != nil {
		t.Fatal("zero CreateInput must leave command settings unspecified")
	}

	explicitlyDisabled := model.CommandsSettings{}
	input.CommandsSettings = &explicitlyDisabled
	if input.CommandsSettings == nil || *input.CommandsSettings != explicitlyDisabled {
		t.Fatal("CreateInput must preserve explicit all-false command settings")
	}

	update := UpdateInput{CommandsSettings: model.CommandsSettings{Mmr: true}}
	if !update.CommandsSettings.Mmr {
		t.Fatal("UpdateInput must retain explicit command settings")
	}
}

func TestCreateUsesEnabledCommandDefaultsWhenSettingsAreUnspecified(t *testing.T) {
	source := readDotaRepositorySource(t, "pgx", "pgx.go")
	want := "COALESCE($8::jsonb, '" + enabledCommandsSettingsJSON + "'::jsonb)"
	if !strings.Contains(source, want) {
		t.Errorf("Dota settings create query must use %q when command settings are nil", want)
	}
}

func TestEnableDotaCommandsByDefaultMigrationContract(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate repository test file")
	}

	matches, err := filepath.Glob(
		filepath.Join(
			filepath.Dir(thisFile),
			"..",
			"..",
			"migrations",
			"postgres",
			"*_enable_dota_commands_by_default.sql",
		),
	)
	if err != nil {
		t.Fatalf("find Dota commands migration: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("found %d Dota commands default migrations, want 1", len(matches))
	}

	contents, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatalf("read Dota commands migration: %v", err)
	}

	for _, fragment := range []string{
		"ALTER COLUMN commands_settings SET DEFAULT '" + enabledCommandsSettingsJSON + "'::jsonb",
		"'" + enabledCommandsSettingsJSON + "'::jsonb || commands_settings",
		"WHERE NOT (commands_settings ?& ARRAY['mmr', 'wl', 'lg', 'gm', 'np', 'wp'])",
		"ALTER COLUMN commands_settings SET DEFAULT '{}'::jsonb",
	} {
		if !strings.Contains(string(contents), fragment) {
			t.Errorf("migration does not contain %q", fragment)
		}
	}
}

func readDotaRepositorySource(t *testing.T, path ...string) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate repository test file")
	}

	contents, err := os.ReadFile(filepath.Join(append([]string{filepath.Dir(thisFile)}, path...)...))
	if err != nil {
		t.Fatalf("read Dota repository source: %v", err)
	}

	return string(contents)
}
