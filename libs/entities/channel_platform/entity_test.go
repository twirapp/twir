package channel_platform

import (
	"encoding/json"
	"testing"
)

func TestChannelPlatformParseTwitchBotConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		botConfig json.RawMessage
		want      TwitchBotConfig
		wantErr   bool
	}{
		{
			name:      "valid full config",
			botConfig: json.RawMessage(`{"bot_id":"bot-1","is_bot_mod":true,"is_twitch_banned":false}`),
			want: TwitchBotConfig{
				BotID:          "bot-1",
				IsBotMod:       true,
				IsTwitchBanned: false,
			},
		},
		{
			name:      "empty raw message",
			botConfig: nil,
			want:      TwitchBotConfig{},
		},
		{
			name:      "empty object",
			botConfig: json.RawMessage(`{}`),
			want:      TwitchBotConfig{},
		},
		{
			name:      "null value",
			botConfig: json.RawMessage(`null`),
			want:      TwitchBotConfig{},
		},
		{
			name:      "unknown fields are ignored",
			botConfig: json.RawMessage(`{"bot_id":"bot-1","unknown_key":"value"}`),
			want:      TwitchBotConfig{BotID: "bot-1"},
		},
		{
			name:      "malformed json",
			botConfig: json.RawMessage(`{`),
			want:      TwitchBotConfig{},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			binding := ChannelPlatform{BotConfig: tt.botConfig}

			got, err := binding.ParseTwitchBotConfig()
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseTwitchBotConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("ParseTwitchBotConfig() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
