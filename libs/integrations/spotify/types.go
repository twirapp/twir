package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type SpotifyTrack struct {
	ID         string `json:"id"`
	URI        string `json:"uri"`
	Name       string `json:"name"`
	ArtistName string `json:"artistName"`
	AlbumName  string `json:"albumName"`
	DurationMs int    `json:"durationMs"`
	ImageURL   string `json:"imageUrl"`
	Type       string `json:"type"`
}

type Device struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	IsActive         bool   `json:"is_active"`
	IsRestricted     bool   `json:"is_restricted"`
	IsPrivateSession bool   `json:"is_private_session"`
}

type CurrentlyPlaying struct {
	CurrentlyPlayingType string        `json:"currently_playing_type"`
	Item                 *SpotifyTrack `json:"item"`
	IsPlaying            bool          `json:"is_playing"`
	ProgressMs           int           `json:"progress_ms"`
	Device               Device        `json:"device"`
}

type spotifyTrackItem struct {
	ID      string `json:"id"`
	URI     string `json:"uri"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
	Album struct {
		Name   string `json:"name"`
		Images []struct {
			URL string `json:"url"`
		} `json:"images"`
	} `json:"album"`
	DurationMs int `json:"duration_ms"`
}

type spotifySearchResponse struct {
	Tracks struct {
		Items []spotifyTrackItem `json:"items"`
	} `json:"tracks"`
}

type spotifyQueueResponse struct {
	Queue []spotifyTrackItem `json:"queue"`
}

type spotifyDevicesResponse struct {
	Devices []Device `json:"devices"`
}

type spotifyCurrentlyPlayingResponse struct {
	CurrentlyPlayingType string            `json:"currently_playing_type"`
	Item                 *spotifyTrackItem `json:"item"`
	IsPlaying            bool              `json:"is_playing"`
	ProgressMs           int               `json:"progress_ms"`
	Device               Device            `json:"device"`
}

type spotifyPlaybackError struct {
	Error struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Reason  string `json:"reason"`
	} `json:"error"`
}

func trackFromItem(item spotifyTrackItem) SpotifyTrack {
	track := SpotifyTrack{
		ID:         item.ID,
		URI:        item.URI,
		Name:       item.Name,
		DurationMs: item.DurationMs,
		Type:       item.Type,
	}

	if len(item.Artists) > 0 {
		track.ArtistName = item.Artists[0].Name
	}
	if item.Album.Name != "" {
		track.AlbumName = item.Album.Name
	}
	if len(item.Album.Images) > 0 {
		track.ImageURL = item.Album.Images[0].URL
	}

	return track
}

func playbackError(statusCode int, body []byte, header http.Header) error {
	switch statusCode {
	case http.StatusForbidden:
		message := playbackErrorMessage(body)
		switch {
		case strings.Contains(message, "premium"):
			return ErrNotPremium
		case strings.Contains(message, "restricted"):
			return ErrRestrictedDevice
		case strings.Contains(message, "scope"):
			return ErrInsufficientScope
		default:
			return ErrInsufficientScope
		}
	case http.StatusNotFound:
		message := playbackErrorMessage(body)
		switch {
		case strings.Contains(message, "active device") || strings.Contains(message, "no device"):
			return ErrNoActiveDevice
		case strings.Contains(message, "not found"):
			return ErrTrackNotFound
		default:
			return ErrNoActiveDevice
		}
	case http.StatusTooManyRequests:
		return rateLimitedError(header)
	default:
		return fmt.Errorf("spotify request failed (status %d): %s", statusCode, truncateBody(body))
	}
}

func truncateBody(body []byte) string {
	const maxLen = 200
	text := strings.TrimSpace(string(body))
	if len(text) > maxLen {
		return text[:maxLen] + "..."
	}
	return text
}

func playbackErrorMessage(body []byte) string {
	if len(body) == 0 {
		return ""
	}

	var parsed spotifyPlaybackError
	if err := json.Unmarshal(body, &parsed); err == nil {
		parts := []string{parsed.Error.Reason, parsed.Error.Message}
		message := strings.ToLower(strings.TrimSpace(strings.Join(parts, " ")))
		if message != "" {
			return message
		}
	}

	return strings.ToLower(strings.TrimSpace(string(body)))
}

func rateLimitedError(header http.Header) error {
	return &RateLimitedError{RetryAfterSeconds: retryAfterSeconds(header)}
}

func retryAfterSeconds(header http.Header) int {
	if header == nil {
		return 0
	}

	value := strings.TrimSpace(header.Get("Retry-After"))
	if value == "" {
		return 0
	}

	seconds, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return seconds
}
