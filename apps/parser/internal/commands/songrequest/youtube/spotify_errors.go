package sr_youtube

import (
	"errors"

	"github.com/twirapp/twir/libs/integrations/spotify"
)

var (
	errSpotifyTrackNotFound           = errors.New("spotify song requests: track not found")
	errSpotifyTrackAlreadyInQueue     = errors.New("spotify song requests: track already in queue")
	errSpotifyMaxRequestsExceeded     = errors.New("spotify song requests: max requests exceeded")
	errSpotifyUserMaxRequestsExceeded = errors.New("spotify song requests: user max requests exceeded")
	errSpotifyDurationNotAllowed      = errors.New("spotify song requests: duration not allowed")
)

func normalizeSpotifySongRequestError(err error) error {
	if err == nil {
		return nil
	}

	switch err.Error() {
	case spotify.ErrNotConnected.Error():
		return spotify.ErrNotConnected
	case spotify.ErrInsufficientScope.Error():
		return spotify.ErrInsufficientScope
	case spotify.ErrNoActiveDevice.Error():
		return spotify.ErrNoActiveDevice
	case errSpotifyTrackNotFound.Error():
		return errSpotifyTrackNotFound
	case errSpotifyTrackAlreadyInQueue.Error():
		return errSpotifyTrackAlreadyInQueue
	case errSpotifyMaxRequestsExceeded.Error():
		return errSpotifyMaxRequestsExceeded
	case errSpotifyUserMaxRequestsExceeded.Error():
		return errSpotifyUserMaxRequestsExceeded
	case errSpotifyDurationNotAllowed.Error():
		return errSpotifyDurationNotAllowed
	default:
		return err
	}
}

func spotifySongRequestErrorMessage(err error) string {
	normalized := normalizeSpotifySongRequestError(err)

	switch {
	case errors.Is(normalized, spotify.ErrNotConnected):
		return "Spotify not connected"
	case errors.Is(normalized, spotify.ErrInsufficientScope):
		return "Spotify needs re-auth (missing playback scope)"
	case errors.Is(normalized, spotify.ErrNoActiveDevice):
		return "No active Spotify device"
	case errors.Is(normalized, errSpotifyTrackNotFound):
		return "Track not found"
	case errors.Is(normalized, errSpotifyTrackAlreadyInQueue):
		return "Song already in queue"
	case errors.Is(normalized, errSpotifyMaxRequestsExceeded):
		return "Queue is full"
	case errors.Is(normalized, errSpotifyUserMaxRequestsExceeded):
		return "You have too many active requests"
	case errors.Is(normalized, errSpotifyDurationNotAllowed):
		return "Track duration not allowed"
	default:
		return ""
	}
}
