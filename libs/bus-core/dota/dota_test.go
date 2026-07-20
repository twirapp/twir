package dota

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"testing"
)

type legacyAegisPickupMessage struct {
	ChannelID    string
	TwitchUserID string
	PlayerName   string
	GameTime     int
}

func TestAegisPickupMessageGobDecodesLegacyMessage(t *testing.T) {
	legacy := legacyAegisPickupMessage{
		ChannelID:    "channel-id",
		TwitchUserID: "twitch-user-id",
		PlayerName:   "Aegis Carrier",
		GameTime:     1_234,
	}

	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(legacy); err != nil {
		t.Fatalf("encode legacy Aegis pickup message: %v", err)
	}

	var decoded AegisPickupMessage
	if err := gob.NewDecoder(&buffer).Decode(&decoded); err != nil {
		t.Fatalf("decode legacy Aegis pickup message: %v", err)
	}

	want := AegisPickupMessage{
		ChannelID:    legacy.ChannelID,
		TwitchUserID: legacy.TwitchUserID,
		PlayerName:   legacy.PlayerName,
		GameTime:     legacy.GameTime,
	}
	if !reflect.DeepEqual(decoded, want) {
		t.Fatalf("decoded legacy Aegis pickup message = %#v, want %#v", decoded, want)
	}
	if decoded.PlayerID != nil {
		t.Fatalf("decoded legacy PlayerID = %v, want nil", *decoded.PlayerID)
	}
}

func TestAegisPickupMessageGobDecodesCurrentMessageIntoLegacyType(t *testing.T) {
	playerID := 42
	current := AegisPickupMessage{
		ChannelID:    "channel-id",
		TwitchUserID: "twitch-user-id",
		PlayerName:   "Aegis Carrier",
		PlayerID:     &playerID,
		GameTime:     1_234,
	}

	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(current); err != nil {
		t.Fatalf("encode current Aegis pickup message: %v", err)
	}

	var decoded legacyAegisPickupMessage
	if err := gob.NewDecoder(&buffer).Decode(&decoded); err != nil {
		t.Fatalf("decode current Aegis pickup message into legacy type: %v", err)
	}

	want := legacyAegisPickupMessage{
		ChannelID:    current.ChannelID,
		TwitchUserID: current.TwitchUserID,
		PlayerName:   current.PlayerName,
		GameTime:     current.GameTime,
	}
	if !reflect.DeepEqual(decoded, want) {
		t.Fatalf("decoded current Aegis pickup message into legacy type = %#v, want %#v", decoded, want)
	}
}

func TestGetDataResponseGobRoundTrip(t *testing.T) {
	want := GetDataResponse{
		Enabled:                 true,
		Linked:                  true,
		InGame:                  true,
		Mmr:                     4_200,
		SessionWins:             3,
		SessionLosses:           1,
		HeroName:                "npc_dota_hero_axe",
		MatchID:                 123_456_789,
		TeamIsRadiant:           true,
		RadiantScore:            24,
		DireScore:               18,
		GameTime:                1_987,
		WinProbability:          0,
		WinProbabilityAvailable: true,
		NotablePlayers:          []string{"Player One", "Player Two"},
		LastGame: &LastGameInfo{
			HeroName:  "npc_dota_hero_juggernaut",
			Kills:     12,
			Deaths:    4,
			Assists:   9,
			Win:       true,
			DurationS: 2_345,
		},
	}

	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(want); err != nil {
		t.Fatalf("encode GetDataResponse: %v", err)
	}

	var decoded GetDataResponse
	if err := gob.NewDecoder(&buffer).Decode(&decoded); err != nil {
		t.Fatalf("decode GetDataResponse: %v", err)
	}

	if !reflect.DeepEqual(decoded, want) {
		t.Fatalf("GetDataResponse Gob round trip = %#v, want %#v", decoded, want)
	}
}
