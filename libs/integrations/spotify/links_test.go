package spotify

import "testing"

func TestParseTrackID(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		wantID string
		wantOk bool
	}{
		{
			name:   "open spotify url with tracking params",
			query:  "https://open.spotify.com/track/6IqfoZlee8TYVSUhIiug0P?si=1a13b96e0ad14bf0",
			wantID: "6IqfoZlee8TYVSUhIiug0P",
			wantOk: true,
		},
		{
			name:   "open spotify url without params",
			query:  "https://open.spotify.com/track/6IqfoZlee8TYVSUhIiug0P",
			wantID: "6IqfoZlee8TYVSUhIiug0P",
			wantOk: true,
		},
		{
			name:   "http url",
			query:  "http://open.spotify.com/track/abc123",
			wantID: "abc123",
			wantOk: true,
		},
		{
			name:   "intl url",
			query:  "https://open.spotify.com/intl-ru/track/abc123?si=x",
			wantID: "abc123",
			wantOk: true,
		},
		{
			name:   "spotify uri",
			query:  "spotify:track:6IqfoZlee8TYVSUhIiug0P",
			wantID: "6IqfoZlee8TYVSUhIiug0P",
			wantOk: true,
		},
		{
			name:   "plain text query",
			query:  "rick astley never gonna give you up",
			wantOk: false,
		},
		{
			name:   "album url is not a track",
			query:  "https://open.spotify.com/album/6IqfoZlee8TYVSUhIiug0P",
			wantOk: false,
		},
		{
			name:   "artist uri is not a track",
			query:  "spotify:artist:6IqfoZlee8TYVSUhIiug0P",
			wantOk: false,
		},
		{
			name:   "whitespace is trimmed",
			query:  "  spotify:track:abc123  ",
			wantID: "abc123",
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, ok := ParseTrackID(tt.query)
			if ok != tt.wantOk {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOk)
			}
			if id != tt.wantID {
				t.Fatalf("id = %q, want %q", id, tt.wantID)
			}
		})
	}
}
