package resolvers

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	now_playing_fetcher "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/now-playing-fetcher"
)

func TestSendNowPlayingTrack(t *testing.T) {
	t.Run("sends value", func(t *testing.T) {
		output := make(chan *gqlmodel.NowPlayingOverlayTrack, 1)
		value := &gqlmodel.NowPlayingOverlayTrack{Artist: "Artist", Title: "Title"}

		if sent := sendNowPlayingTrack(context.Background(), output, value); !sent {
			t.Fatal("expected value to be sent")
		}

		if got := <-output; got != value {
			t.Fatalf("expected track pointer %p, got %p", value, got)
		}
	})

	t.Run("cancellation unblocks send with no receiver", func(t *testing.T) {
		baseCtx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		ctx := &sendReadyContext{
			Context: baseCtx,
			ready:   make(chan struct{}),
		}
		output := make(chan *gqlmodel.NowPlayingOverlayTrack)
		result := make(chan bool, 1)

		go func() {
			result <- sendNowPlayingTrack(
				ctx,
				output,
				&gqlmodel.NowPlayingOverlayTrack{Artist: "Artist", Title: "Title"},
			)
		}()

		waitForStreamSignal(t, ctx.ready)
		select {
		case sent := <-result:
			t.Fatalf("send returned before cancellation: sent=%t", sent)
		default:
		}

		cancel()
		select {
		case sent := <-result:
			if sent {
				t.Fatal("expected canceled send to return false")
			}
		case <-time.After(streamTestGuard):
			t.Fatal("timed out waiting for canceled send")
		}
	})
}

const (
	streamTestGuard  = 500 * time.Millisecond
	streamRetryGuard = 2 * time.Second
)

func TestStreamNowPlaying(t *testing.T) {
	t.Run("cancellation interrupts interval after an emission", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		output := make(chan *gqlmodel.NowPlayingOverlayTrack)
		done := make(chan struct{})

		go func() {
			streamNowPlaying(
				ctx,
				output,
				func(context.Context) (*now_playing_fetcher.Track, error) {
					return &now_playing_fetcher.Track{Artist: "Artist", Title: "Title"}, nil
				},
				func(error) {},
				time.Second,
			)
			close(done)
		}()

		select {
		case track := <-output:
			if track == nil || track.Title != "Title" {
				t.Fatalf("expected emitted track, got %#v", track)
			}
		case <-time.After(streamTestGuard):
			t.Fatal("timed out waiting for track emission")
		}

		cancel()
		waitForStreamSignal(t, done)

		if _, ok := <-output; ok {
			t.Fatal("expected output channel to be closed")
		}
	})

	t.Run("emits nil track", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		output := make(chan *gqlmodel.NowPlayingOverlayTrack)
		done := make(chan struct{})

		go func() {
			streamNowPlaying(
				ctx,
				output,
				func(context.Context) (*now_playing_fetcher.Track, error) { return nil, nil },
				func(error) {},
				time.Second,
			)
			close(done)
		}()

		select {
		case track := <-output:
			if track != nil {
				t.Fatalf("expected nil track, got %#v", track)
			}
		case <-time.After(streamTestGuard):
			t.Fatal("timed out waiting for nil emission")
		}

		cancel()
		waitForStreamSignal(t, done)
	})

	t.Run("reports fetch errors retries and cancels during interval", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		output := make(chan *gqlmodel.NowPlayingOverlayTrack)
		attempts := make(chan int, 2)
		reported := make(chan error, 2)
		done := make(chan struct{})
		fetchErr := errors.New("fetch failed")
		attempt := 0

		go func() {
			streamNowPlaying(
				ctx,
				output,
				func(context.Context) (*now_playing_fetcher.Track, error) {
					attempt++
					attempts <- attempt
					return nil, fetchErr
				},
				func(err error) { reported <- err },
				time.Second,
			)
			close(done)
		}()

		for expectedAttempt := 1; expectedAttempt <= 2; expectedAttempt++ {
			select {
			case got := <-attempts:
				if got != expectedAttempt {
					t.Fatalf("expected fetch attempt %d, got %d", expectedAttempt, got)
				}
			case <-time.After(streamRetryGuard):
				t.Fatalf("timed out waiting for fetch attempt %d", expectedAttempt)
			}

			select {
			case got := <-reported:
				if !errors.Is(got, fetchErr) {
					t.Fatalf("expected fetch error %v, got %v", fetchErr, got)
				}
			case <-time.After(streamRetryGuard):
				t.Fatalf("timed out waiting for error callback %d", expectedAttempt)
			}
		}

		cancel()
		waitForStreamSignal(t, done)

		if _, ok := <-output; ok {
			t.Fatal("expected output channel to be closed")
		}
	})
}

func TestMapNowPlayingTrack(t *testing.T) {
	t.Run("maps timed track with image", func(t *testing.T) {
		progressMs := 12_345
		durationMs := 234_567
		track := &now_playing_fetcher.Track{
			Artist:     "Artist",
			Title:      "Title",
			ImageUrl:   "https://example.com/cover.jpg",
			ProgressMs: &progressMs,
			DurationMs: &durationMs,
		}
		want := &gqlmodel.NowPlayingOverlayTrack{
			Artist:     track.Artist,
			Title:      track.Title,
			ImageURL:   &track.ImageUrl,
			ProgressMs: track.ProgressMs,
			DurationMs: track.DurationMs,
		}

		if got := mapNowPlayingTrack(track); !reflect.DeepEqual(got, want) {
			t.Fatalf("expected %#v, got %#v", want, got)
		}
	})

	t.Run("maps ambient track without timing", func(t *testing.T) {
		track := &now_playing_fetcher.Track{Artist: "Ambient Artist", Title: "Ambient Title"}
		want := &gqlmodel.NowPlayingOverlayTrack{Artist: track.Artist, Title: track.Title}

		if got := mapNowPlayingTrack(track); !reflect.DeepEqual(got, want) {
			t.Fatalf("expected %#v, got %#v", want, got)
		}
	})

	t.Run("leaves empty image nil", func(t *testing.T) {
		got := mapNowPlayingTrack(&now_playing_fetcher.Track{})
		if got == nil {
			t.Fatal("expected mapped track, got nil")
		}
		if got.ImageURL != nil {
			t.Fatalf("expected nil image URL, got %q", *got.ImageURL)
		}
	})

	t.Run("maps nil input to nil", func(t *testing.T) {
		if got := mapNowPlayingTrack(nil); got != nil {
			t.Fatalf("expected nil track, got %#v", got)
		}
	})
}

func waitForStreamSignal(t *testing.T, signal <-chan struct{}) {
	t.Helper()

	select {
	case <-signal:
	case <-time.After(streamTestGuard):
		t.Fatal("timed out waiting for stream helper")
	}
}

type sendReadyContext struct {
	context.Context
	ready chan struct{}
	once  sync.Once
}

func (c *sendReadyContext) Done() <-chan struct{} {
	c.once.Do(func() { close(c.ready) })
	return c.Context.Done()
}
