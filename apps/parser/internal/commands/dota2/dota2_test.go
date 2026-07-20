package dota2

import (
	"testing"

	busdota "github.com/twirapp/twir/libs/bus-core/dota"
)

func TestMedalTierBoundaries(t *testing.T) {
	tests := []struct {
		name string
		mmr  int
		want medalTier
	}{
		{name: "herald starts at zero", mmr: 0, want: medalHerald},
		{name: "herald ends at 769", mmr: 769, want: medalHerald},
		{name: "guardian starts at 770", mmr: 770, want: medalGuardian},
		{name: "guardian ends at 1539", mmr: 1539, want: medalGuardian},
		{name: "crusader starts at 1540", mmr: 1540, want: medalCrusader},
		{name: "crusader ends at 2309", mmr: 2309, want: medalCrusader},
		{name: "archon starts at 2310", mmr: 2310, want: medalArchon},
		{name: "archon ends at 3079", mmr: 3079, want: medalArchon},
		{name: "legend starts at 3080", mmr: 3080, want: medalLegend},
		{name: "legend ends at 3849", mmr: 3849, want: medalLegend},
		{name: "ancient starts at 3850", mmr: 3850, want: medalAncient},
		{name: "ancient ends at 4619", mmr: 4619, want: medalAncient},
		{name: "divine starts at 4620", mmr: 4620, want: medalDivine},
		{name: "divine ends at 5419", mmr: 5419, want: medalDivine},
		{name: "immortal starts at 5420", mmr: 5420, want: medalImmortal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := medalForMMR(tt.mmr); got != tt.want {
				t.Fatalf("medalForMMR(%d) = %q, want %q", tt.mmr, got, tt.want)
			}
		})
	}
}

func TestFormatWinLoss(t *testing.T) {
	tests := []struct {
		wins   int
		losses int
		want   winLossOutput
	}{
		{wins: 0, losses: 0, want: winLossOutput{Record: "0-0", WinRate: "0.0%"}},
		{wins: 1, losses: 2, want: winLossOutput{Record: "1-2", WinRate: "33.3%"}},
		{wins: 3, losses: 1, want: winLossOutput{Record: "3-1", WinRate: "75.0%"}},
	}

	for _, tt := range tests {
		if got := formatWinLoss(tt.wins, tt.losses); got != tt.want {
			t.Fatalf("formatWinLoss(%d, %d) = %#v, want %#v", tt.wins, tt.losses, got, tt.want)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		seconds int
		want    string
	}{
		{seconds: 0, want: "0:00"},
		{seconds: 65, want: "1:05"},
		{seconds: 3661, want: "61:01"},
	}

	for _, tt := range tests {
		if got := formatDuration(tt.seconds); got != tt.want {
			t.Fatalf("formatDuration(%d) = %q, want %q", tt.seconds, got, tt.want)
		}
	}
}

func TestFormatLastGame(t *testing.T) {
	if _, ok := formatLastGame(nil); ok {
		t.Fatal("formatLastGame(nil) reported a game")
	}

	got, ok := formatLastGame(&busdota.LastGameInfo{
		HeroName:  "Axe",
		Kills:     12,
		Deaths:    3,
		Assists:   7,
		Win:       true,
		DurationS: 125,
	})
	if !ok {
		t.Fatal("formatLastGame reported no game")
	}

	want := lastGameOutput{
		HeroName: "Axe",
		KDA:      "12/3/7",
		Won:      true,
		Duration: "2:05",
	}
	if got != want {
		t.Fatalf("formatLastGame() = %#v, want %#v", got, want)
	}
}

func TestFormatWinProbability(t *testing.T) {
	tests := []struct {
		probability float64
		want        string
	}{
		{probability: 0, want: "0.0%"},
		{probability: 0.625, want: "62.5%"},
		{probability: 1, want: "100.0%"},
	}

	for _, tt := range tests {
		if got := formatWinProbability(tt.probability); got != tt.want {
			t.Fatalf("formatWinProbability(%f) = %q, want %q", tt.probability, got, tt.want)
		}
	}
}

func TestWinProbabilityOutputDistinguishesUnavailableFromZero(t *testing.T) {
	output, available := winProbabilityOutput(&busdota.GetDataResponse{})
	if available {
		t.Error("unavailable win probability was reported as available")
	}
	if output != "" {
		t.Errorf("unavailable win probability output = %q, want empty", output)
	}

	output, available = winProbabilityOutput(&busdota.GetDataResponse{
		WinProbabilityAvailable: true,
		WinProbability:          0,
	})
	if !available {
		t.Error("valid zero win probability was reported as unavailable")
	}
	if output != "0.0%" {
		t.Errorf("zero win probability output = %q, want 0.0%%", output)
	}
}

func TestJoinNotablePlayers(t *testing.T) {
	tests := []struct {
		players []string
		want    string
	}{
		{players: nil, want: ""},
		{players: []string{"Player One"}, want: "Player One"},
		{players: []string{"Player One", "Player Two"}, want: "Player One, Player Two"},
	}

	for _, tt := range tests {
		if got := joinNotablePlayers(tt.players); got != tt.want {
			t.Fatalf("joinNotablePlayers(%q) = %q, want %q", tt.players, got, tt.want)
		}
	}
}
