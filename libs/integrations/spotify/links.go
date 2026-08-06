package spotify

import (
	"regexp"
	"strings"
)

var (
	spotifyTrackURLPattern = regexp.MustCompile(
		`^https?://open\.spotify\.com/(?:intl-[a-z]{2}(?:-[a-z]{2})?/)?track/([A-Za-z0-9]+)`,
	)
	spotifyTrackURIPattern = regexp.MustCompile(`^spotify:track:([A-Za-z0-9]+)$`)
)

// ParseTrackID extracts a track id from a Spotify link or URI.
func ParseTrackID(query string) (string, bool) {
	query = strings.TrimSpace(query)
	if match := spotifyTrackURIPattern.FindStringSubmatch(query); match != nil {
		return match[1], true
	}
	if match := spotifyTrackURLPattern.FindStringSubmatch(query); match != nil {
		return match[1], true
	}
	return "", false
}
