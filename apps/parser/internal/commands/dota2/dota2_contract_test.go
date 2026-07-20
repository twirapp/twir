package dota2

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/twirapp/twir/apps/parser/internal/types"
)

func TestReadCommandsAreVisible(t *testing.T) {
	commands := map[string]*types.DefaultCommand{
		"mmr": Mmr,
		"wl":  Wl,
		"lg":  Lg,
		"gm":  Gm,
		"np":  Np,
		"wp":  Wp,
	}

	for name, command := range commands {
		if !command.ChannelsCommands.Visible {
			t.Errorf("%s command is not publicly visible", name)
		}
	}
}

func TestMmrSetUsesUpdateToTouchUpdatedAt(t *testing.T) {
	source := readDota2Source(t, "mmr.go")
	if strings.Contains(source, `UpdateColumn("mmr", mmr)`) {
		t.Error("MmrSet must not use UpdateColumn because it bypasses updated_at")
	}
	if !strings.Contains(source, `Update("mmr", mmr)`) {
		t.Error("MmrSet must use Update so GORM updates updated_at")
	}
}

func TestGetDotaDataUsesChatSafeDeadline(t *testing.T) {
	source := readDota2Source(t, "dota2.go")
	for _, fragment := range []string{
		"dotaDataRequestTimeout = 5 * time.Second",
		"context.WithTimeout(ctx, dotaDataRequestTimeout)",
		"defer cancel()",
		"Request(\n\t\trequestCtx,",
	} {
		if !strings.Contains(source, fragment) {
			t.Errorf("getDotaData must contain %q", fragment)
		}
	}
}

func TestWpRejectsUnavailableWinProbability(t *testing.T) {
	source := readDota2Source(t, "wp.go")
	for _, fragment := range []string{
		"winProbabilityOutput(data)",
		"if !available",
		"Errors.WinProbabilityUnavailable",
	} {
		if !strings.Contains(source, fragment) {
			t.Errorf("Wp handler must contain %q", fragment)
		}
	}
}

func readDota2Source(t *testing.T, name string) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate Dota command test file")
	}

	contents, err := os.ReadFile(filepath.Join(filepath.Dir(thisFile), name))
	if err != nil {
		t.Fatalf("read Dota command source: %v", err)
	}

	return string(contents)
}
