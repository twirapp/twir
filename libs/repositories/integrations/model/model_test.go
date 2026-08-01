package model

import "testing"

func TestServiceStreamElementsValue(t *testing.T) {
	t.Parallel()
	if got := ServiceStreamElements; got != Service("STREAMELEMENTS") {
		t.Fatalf("ServiceStreamElements = %q", got)
	}
}
