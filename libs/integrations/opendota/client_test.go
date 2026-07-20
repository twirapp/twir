package opendota

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeroes_DecodesIDKeyedResponse(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/constants/heroes" {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(
				w,
				`{"1":{"id":999,"name":"npc_dota_hero_antimage","localized_name":"Anti-Mage"},"2":{"id":2,"name":"npc_dota_hero_axe","localized_name":"Axe"}}`,
			)
		}),
	)
	t.Cleanup(server.Close)

	heroes, err := New(WithBaseURL(server.URL)).Heroes(context.Background())
	if err != nil {
		t.Fatalf("Heroes returned error: %v", err)
	}

	if got := heroes[1]; got != "Anti-Mage" {
		t.Errorf("expected hero 1 to be Anti-Mage, got %q", got)
	}
	if got := heroes[2]; got != "Axe" {
		t.Errorf("expected hero 2 to be Axe, got %q", got)
	}
	if _, ok := heroes[999]; ok {
		t.Error("expected response object key to identify the hero, not payload id")
	}
}
