package opendota

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestPlayerHeroes_DecodesResponse(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet || r.URL.Path != "/players/123456/heroes" {
				http.NotFound(w, r)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `[{"hero_id":1,"games":32,"win":19},{"hero_id":2,"games":7,"win":2}]`)
		}),
	)
	t.Cleanup(server.Close)

	heroes, err := New(WithBaseURL(server.URL)).PlayerHeroes(context.Background(), 123456)
	if err != nil {
		t.Fatalf("PlayerHeroes returned error: %v", err)
	}
	if len(heroes) != 2 {
		t.Fatalf("expected 2 player heroes, got %v", heroes)
	}
	if heroes[0] != (PlayerHero{HeroID: 1, Games: 32, Win: 19}) {
		t.Errorf("unexpected first player hero: %+v", heroes[0])
	}
	if heroes[1] != (PlayerHero{HeroID: 2, Games: 7, Win: 2}) {
		t.Errorf("unexpected second player hero: %+v", heroes[1])
	}
}

func TestProPlayers_RejectsOversizedResponseBody(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprint(w, "["+strings.Repeat(" ", 1<<20)+"]")
		}),
	)
	t.Cleanup(server.Close)

	_, err := New(WithBaseURL(server.URL)).ProPlayers(context.Background())
	if err == nil {
		t.Fatal("expected oversized response body error")
	}
	if !strings.Contains(err.Error(), "response body exceeds") {
		t.Errorf("expected response body limit error, got %v", err)
	}
}

func TestProPlayers_TruncatesNonSuccessResponsePreview(t *testing.T) {
	body := strings.Repeat("x", 8<<10)
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, body, http.StatusInternalServerError)
		}),
	)
	t.Cleanup(server.Close)

	_, err := New(WithBaseURL(server.URL)).ProPlayers(context.Background())
	if err == nil {
		t.Fatal("expected unexpected status error")
	}
	if !strings.HasSuffix(err.Error(), "...") {
		t.Errorf("expected truncated response preview, got %q", err)
	}
	if strings.Contains(err.Error(), body) {
		t.Error("error included the full response body")
	}
}
