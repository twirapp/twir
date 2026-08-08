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

type legacyMatchStartedMessage struct {
	ChannelID      string
	TwitchUserID   string
	SteamAccountID string
	HeroName       string
}

type legacyMatchEndedMessage struct {
	ChannelID      string
	TwitchUserID   string
	SteamAccountID string
	Win            bool
	HeroName       string
	Mmr            int
	SessionWins    int
	SessionLosses  int
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

func TestMatchStartedMessageGobDecodesLegacyMessage(t *testing.T) {
	legacy := legacyMatchStartedMessage{
		ChannelID:      "channel-id",
		TwitchUserID:   "twitch-user-id",
		SteamAccountID: "steam-account-id",
		HeroName:       "axe",
	}

	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(legacy); err != nil {
		t.Fatalf("encode legacy match started message: %v", err)
	}

	var decoded MatchStartedMessage
	if err := gob.NewDecoder(&buffer).Decode(&decoded); err != nil {
		t.Fatalf("decode legacy match started message: %v", err)
	}

	want := MatchStartedMessage{
		ChannelID:      legacy.ChannelID,
		TwitchUserID:   legacy.TwitchUserID,
		SteamAccountID: legacy.SteamAccountID,
		HeroName:       legacy.HeroName,
	}
	if !reflect.DeepEqual(decoded, want) {
		t.Fatalf("decoded legacy match started message = %#v, want %#v", decoded, want)
	}
}

func TestMatchStartedMessageGobDecodesCurrentMessageIntoLegacyType(t *testing.T) {
	current := MatchStartedMessage{
		ChannelID:      "channel-id",
		TwitchUserID:   "twitch-user-id",
		SteamAccountID: "steam-account-id",
		HeroName:       "axe",
		MatchID:        123_456,
		TeamKnown:      true,
	}

	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(current); err != nil {
		t.Fatalf("encode current match started message: %v", err)
	}

	var decoded legacyMatchStartedMessage
	if err := gob.NewDecoder(&buffer).Decode(&decoded); err != nil {
		t.Fatalf("decode current match started message into legacy type: %v", err)
	}

	want := legacyMatchStartedMessage{
		ChannelID:      current.ChannelID,
		TwitchUserID:   current.TwitchUserID,
		SteamAccountID: current.SteamAccountID,
		HeroName:       current.HeroName,
	}
	if !reflect.DeepEqual(decoded, want) {
		t.Fatalf("decoded current match started message into legacy type = %#v, want %#v", decoded, want)
	}
}

func TestMatchEndedMessageGobDecodesLegacyMessage(t *testing.T) {
	legacy := legacyMatchEndedMessage{
		ChannelID:      "channel-id",
		TwitchUserID:   "twitch-user-id",
		SteamAccountID: "steam-account-id",
		Win:            true,
		HeroName:       "axe",
		Mmr:            4_200,
		SessionWins:    12,
		SessionLosses:  3,
	}

	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(legacy); err != nil {
		t.Fatalf("encode legacy match ended message: %v", err)
	}

	var decoded MatchEndedMessage
	if err := gob.NewDecoder(&buffer).Decode(&decoded); err != nil {
		t.Fatalf("decode legacy match ended message: %v", err)
	}

	want := MatchEndedMessage{
		ChannelID:      legacy.ChannelID,
		TwitchUserID:   legacy.TwitchUserID,
		SteamAccountID: legacy.SteamAccountID,
		Win:            legacy.Win,
		HeroName:       legacy.HeroName,
		Mmr:            legacy.Mmr,
		SessionWins:    legacy.SessionWins,
		SessionLosses:  legacy.SessionLosses,
	}
	if !reflect.DeepEqual(decoded, want) {
		t.Fatalf("decoded legacy match ended message = %#v, want %#v", decoded, want)
	}
}

func TestMatchEndedMessageGobDecodesCurrentMessageIntoLegacyType(t *testing.T) {
	current := MatchEndedMessage{
		ChannelID:      "channel-id",
		TwitchUserID:   "twitch-user-id",
		SteamAccountID: "steam-account-id",
		MatchID:        123_456,
		Win:            true,
		HeroName:       "axe",
		Mmr:            4_200,
		SessionWins:    12,
		SessionLosses:  3,
	}

	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(current); err != nil {
		t.Fatalf("encode current match ended message: %v", err)
	}

	var decoded legacyMatchEndedMessage
	if err := gob.NewDecoder(&buffer).Decode(&decoded); err != nil {
		t.Fatalf("decode current match ended message into legacy type: %v", err)
	}

	want := legacyMatchEndedMessage{
		ChannelID:      current.ChannelID,
		TwitchUserID:   current.TwitchUserID,
		SteamAccountID: current.SteamAccountID,
		Win:            current.Win,
		HeroName:       current.HeroName,
		Mmr:            current.Mmr,
		SessionWins:    current.SessionWins,
		SessionLosses:  current.SessionLosses,
	}
	if !reflect.DeepEqual(decoded, want) {
		t.Fatalf("decoded current match ended message into legacy type = %#v, want %#v", decoded, want)
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
