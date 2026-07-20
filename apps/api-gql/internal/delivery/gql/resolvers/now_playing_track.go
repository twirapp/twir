package resolvers

import (
	"context"
	"time"

	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	now_playing_fetcher "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/now-playing-fetcher"
)

func streamNowPlaying(
	ctx context.Context,
	output chan<- *gqlmodel.NowPlayingOverlayTrack,
	fetch func(context.Context) (*now_playing_fetcher.Track, error),
	onError func(error),
	interval time.Duration,
) {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		track, err := fetch(ctx)
		select {
		case <-ctx.Done():
			return
		default:
		}

		if err != nil {
			onError(err)
		} else if !sendNowPlayingTrack(ctx, output, mapNowPlayingTrack(track)) {
			return
		}

		timer := time.NewTimer(interval)
		select {
		case <-timer.C:
		case <-ctx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			return
		}
	}
}

func sendNowPlayingTrack(
	ctx context.Context,
	output chan<- *gqlmodel.NowPlayingOverlayTrack,
	value *gqlmodel.NowPlayingOverlayTrack,
) bool {
	select {
	case output <- value:
		return true
	case <-ctx.Done():
		return false
	}
}

func mapNowPlayingTrack(track *now_playing_fetcher.Track) *gqlmodel.NowPlayingOverlayTrack {
	if track == nil {
		return nil
	}

	var imageURL *string
	if track.ImageUrl != "" {
		imageURL = &track.ImageUrl
	}

	return &gqlmodel.NowPlayingOverlayTrack{
		Artist:     track.Artist,
		Title:      track.Title,
		ImageURL:   imageURL,
		ProgressMs: track.ProgressMs,
		DurationMs: track.DurationMs,
	}
}
