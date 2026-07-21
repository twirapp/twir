package channel

import (
	"testing"

	"github.com/twirapp/twir/libs/entities/platform"
)

func TestPlatformsForRequest(t *testing.T) {
	t.Parallel()

	t.Run("uses all platforms when none are requested", func(t *testing.T) {
		t.Parallel()

		got := platformsForRequest(nil)
		want := platform.All()

		if len(got) != len(want) {
			t.Fatalf("expected %d platforms, got %d", len(want), len(got))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("expected platform %q at index %d, got %q", want[i], i, got[i])
			}
		}
	})

	t.Run("uses explicitly requested platforms", func(t *testing.T) {
		t.Parallel()

		got := platformsForRequest([]platform.Platform{platform.PlatformKick})

		if len(got) != 1 || got[0] != platform.PlatformKick {
			t.Fatalf("expected only Kick, got %#v", got)
		}
	})
}
