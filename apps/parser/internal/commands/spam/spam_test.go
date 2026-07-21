package spam

import "testing"

func TestStripRepeatVariable(t *testing.T) {
	got := stripRepeatVariable("Join $(repeat) https://t.me/example")
	want := "Join  https://t.me/example"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
