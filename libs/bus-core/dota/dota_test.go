package dota

import (
	"bytes"
	"encoding/gob"
	"testing"
)

type legacyAegisPickupMessage struct {
	ChannelID    string
	TwitchUserID string
	PlayerName   string
	GameTime     int
}

func TestAegisPickupMessageGobCompatibility(t *testing.T) {
	t.Run("legacy payload decodes without player ID", func(t *testing.T) {
		legacy := legacyAegisPickupMessage{
			ChannelID:    "channel",
			TwitchUserID: "twitch-user",
			PlayerName:   "Puppey",
			GameTime:     600,
		}

		var payload bytes.Buffer
		if err := gob.NewEncoder(&payload).Encode(legacy); err != nil {
			t.Fatalf("encode legacy payload: %v", err)
		}

		var current AegisPickupMessage
		if err := gob.NewDecoder(&payload).Decode(&current); err != nil {
			t.Fatalf("decode legacy payload: %v", err)
		}

		if current.ChannelID != legacy.ChannelID ||
			current.TwitchUserID != legacy.TwitchUserID ||
			current.PlayerName != legacy.PlayerName ||
			current.GameTime != legacy.GameTime {
			t.Fatalf("decoded legacy payload differs: %#v", current)
		}
		if current.PlayerID != nil {
			t.Fatalf("expected nil player ID, got %d", *current.PlayerID)
		}
	})

	t.Run("current payload decodes in legacy receiver", func(t *testing.T) {
		playerID := 2
		current := AegisPickupMessage{
			ChannelID:    "channel",
			TwitchUserID: "twitch-user",
			PlayerName:   "Puppey",
			PlayerID:     &playerID,
			GameTime:     600,
		}

		var payload bytes.Buffer
		if err := gob.NewEncoder(&payload).Encode(current); err != nil {
			t.Fatalf("encode current payload: %v", err)
		}

		var legacy legacyAegisPickupMessage
		if err := gob.NewDecoder(&payload).Decode(&legacy); err != nil {
			t.Fatalf("decode current payload: %v", err)
		}

		if legacy.ChannelID != current.ChannelID ||
			legacy.TwitchUserID != current.TwitchUserID ||
			legacy.PlayerName != current.PlayerName ||
			legacy.GameTime != current.GameTime {
			t.Fatalf("decoded current payload differs: %#v", legacy)
		}
	})
}
