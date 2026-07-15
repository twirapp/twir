# Now Playing Overlay Presets Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add four artwork-colored now-playing presets with real Spotify timing and animated ambient fallbacks for Last.fm and VK.

**Architecture:** Extend the Spotify integration and GraphQL subscription with nullable timing, compensate cached progress in the fetcher, and keep rendering logic in `@twir/frontend-now-playing`. Pure timing and palette helpers receive unit tests; Vue composables own lifecycle, and four focused preset components consume the shared state in both dashboard preview and OBS.

**Tech Stack:** Go 1.26, gqlgen GraphQL, Vue 3 `<script setup>`, TypeScript, Bun test runner, CSS container queries and animations, Nuxt dashboard, Vite overlays app.

**Commit policy:** Commit commands below are checkpoints required by the workflow. Run them only if the user explicitly authorizes commits; otherwise skip the commit step and continue verification without staging files.

---

## File Map

Backend timing:

- Modify `libs/integrations/spotify/track.go`: decode and return Spotify progress and duration.
- Create `libs/integrations/spotify/track_test.go`: cover both Spotify endpoint paths.
- Modify `apps/api-gql/internal/delivery/gql/now-playing-fetcher/now-playing-fetcher.go`: carry nullable timing and compensate Redis-cached progress.
- Create `apps/api-gql/internal/delivery/gql/now-playing-fetcher/now-playing-fetcher_test.go`: cover progress compensation and clamping.

GraphQL and enums:

- Modify `apps/api-gql/internal/delivery/gql/schema/overlays/overlays.graphql`: add four preset values and nullable timing fields.
- Modify `libs/types/types/api/overlays/nowplaying.go`: add the persisted Go enum values.
- Regenerate `libs/types/src/api.ts`, `apps/api-gql/internal/delivery/gql/graph/generated.go`, and `apps/api-gql/internal/delivery/gql/gqlmodel/gqlmodel.go`.
- Create `apps/api-gql/internal/delivery/gql/resolvers/now_playing_track.go`: isolate internal-track to GraphQL mapping.
- Create `apps/api-gql/internal/delivery/gql/resolvers/now_playing_track_test.go`: verify nullable timing mapping.
- Modify `apps/api-gql/internal/delivery/gql/resolvers/overlays.resolver.service.go`: use the mapper.

Shared frontend behavior:

- Modify `libs/frontend-now-playing/src/types.ts`: add timing fields and preset names.
- Create `libs/frontend-now-playing/src/utils/progress.ts`: pure timing normalization and formatting.
- Create `libs/frontend-now-playing/tests/progress.test.ts`: timing unit tests.
- Create `libs/frontend-now-playing/src/composables/use-track-progress.ts`: local timing interpolation.
- Create `libs/frontend-now-playing/src/utils/palette.ts`: pure sampling, color derivation, and contrast helpers.
- Create `libs/frontend-now-playing/tests/palette.test.ts`: palette unit tests.
- Create `libs/frontend-now-playing/src/composables/use-artwork-palette.ts`: image/canvas lifecycle and fallback.
- Modify `libs/frontend-now-playing/src/assets/style.css`: shared progress, marquee, ambient, reduced-motion, and palette transitions.

Preset rendering:

- Create `libs/frontend-now-playing/src/presets/pulse-strip.vue`.
- Create `libs/frontend-now-playing/src/presets/aura-stack.vue`.
- Create `libs/frontend-now-playing/src/presets/vinyl-haze.vue`.
- Create `libs/frontend-now-playing/src/presets/signal-deck.vue`.
- Modify `libs/frontend-now-playing/src/now-playing.vue`: select the new components.

Consumers:

- Modify `frontend/overlays/src/composables/now-playing/use-now-playing-socket.ts`: subscribe to timing.
- Regenerate `frontend/overlays/src/gql/graphql.ts` and related client-preset files.
- Modify `web/layers/dashboard/pages/dashboard/overlays/now-playing.vue`: subscribe to timing and provide preview timing.
- Modify `web/layers/dashboard/pages/dashboard/overlays/now-playing/now-playing-form.vue`: expose the four labels.
- Regenerate `web/app/gql/graphql.ts`, `web/app/gql/gql.ts`, and related client-preset files.

## Task 1: Decode Spotify Timing

**Files:**

- Create: `libs/integrations/spotify/track_test.go`
- Modify: `libs/integrations/spotify/track.go`

- [ ] **Step 1: Write endpoint-level failing tests**

Create tests that replace `http.DefaultClient` with a deterministic transport and call the public `GetTrack` method through both scope paths:

```go
package spotify

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func installSpotifyResponse(t *testing.T, body string) {
	t.Helper()
	oldClient := http.DefaultClient
	http.DefaultClient = &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})}
	t.Cleanup(func() { http.DefaultClient = oldClient })
}

func TestGetTrackIncludesTiming(t *testing.T) {
	tests := []struct {
		name   string
		scopes []string
		body   string
	}{
		{
			name: "currently playing endpoint",
			body: `{"progress_ms":70000,"is_playing":true,"item":{"name":"Heat Waves","duration_ms":238000,"artists":[{"name":"Glass Animals"}],"album":{"images":[{"url":"cover"}]}}}`,
		},
		{
			name:   "player state endpoint",
			scopes: []string{"user-read-playback-state"},
			body:   `{"progress_ms":70000,"is_playing":true,"context":{"type":""},"item":{"name":"Heat Waves","duration_ms":238000,"artists":[{"name":"Glass Animals"}],"album":{"images":[{"url":"cover"}]}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installSpotifyResponse(t, tt.body)
			track, err := NewStatic("token", tt.scopes).GetTrack(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			if track.ProgressMs != 70000 || track.DurationMs != 238000 {
				t.Fatalf("unexpected timing: progress=%d duration=%d", track.ProgressMs, track.DurationMs)
			}
		})
	}
}
```

- [ ] **Step 2: Run the focused test and verify failure**

Run from `libs/integrations`:

```bash
go test ./spotify -run TestGetTrackIncludesTiming -v
```

Expected: compilation fails because `GetTrackResponse` has no `ProgressMs` or `DurationMs`.

- [ ] **Step 3: Add timing fields to both decoded response shapes**

In `track.go`, add root progress and item duration fields:

```go
type spotifyCurrentPlayingTrack struct {
	Artists    []spotifyCurrentPlayingArtist `json:"artists"`
	Name       string                        `json:"name"`
	Album      spotifyCurrentPlayingAlbum    `json:"album"`
	DurationMs int                           `json:"duration_ms"`
}

type spotifyCurrentPlayingResponse struct {
	Track      *spotifyCurrentPlayingTrack `json:"item"`
	ProgressMs int                         `json:"progress_ms"`
	IsPlaying  bool                        `json:"is_playing"`
}

type GetTrackResponse struct {
	Playlist   *GetTrackResponsePlaylist `json:"playlist"`
	Title      string                    `json:"title"`
	Artist     string                    `json:"artist"`
	Image      string                    `json:"image"`
	ProgressMs int                       `json:"progressMs"`
	DurationMs int                       `json:"durationMs"`
	IsPlaying  bool                      `json:"isPlaying"`
}
```

Add `ProgressMs int ` + "`json:\"progress_ms\"`" + ` to `spotifyPlayerStateResponse`, then map `data.ProgressMs` and the corresponding item duration into both returned `GetTrackResponse` values.

- [ ] **Step 4: Run and format**

Run from `libs/integrations`:

```bash
gofmt -w spotify/track.go spotify/track_test.go
```

Expected: both subtests pass.

- [ ] **Step 5: Commit checkpoint, only with authorization**

```bash
git add libs/integrations/spotify/track.go libs/integrations/spotify/track_test.go
git commit -m "feat(spotify): expose playback timing"
```

## Task 2: Compensate Cached Spotify Progress

**Files:**

- Create: `apps/api-gql/internal/delivery/gql/now-playing-fetcher/now-playing-fetcher_test.go`
- Modify: `apps/api-gql/internal/delivery/gql/now-playing-fetcher/now-playing-fetcher.go`

- [ ] **Step 1: Write failing tests for cached progress**

```go
package now_playing_fetcher

import (
	"testing"
	"time"
)

func intPointer(value int) *int { return &value }

func TestTrackAdvanceProgress(t *testing.T) {
	now := time.Date(2026, time.July, 12, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		track    Track
		expected int
	}{
		{
			name: "adds elapsed cache time",
			track: Track{ProgressMs: intPointer(10000), DurationMs: intPointer(30000), ProgressObservedAt: now.Add(-5 * time.Second)},
			expected: 15000,
		},
		{
			name: "clamps to duration",
			track: Track{ProgressMs: intPointer(29000), DurationMs: intPointer(30000), ProgressObservedAt: now.Add(-5 * time.Second)},
			expected: 30000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.track.advanceProgress(now)
			if tt.track.ProgressMs == nil || *tt.track.ProgressMs != tt.expected {
				t.Fatalf("expected %d, got %v", tt.expected, tt.track.ProgressMs)
			}
			if !tt.track.ProgressObservedAt.Equal(now) {
				t.Fatalf("observation time was not refreshed")
			}
		})
	}
}

func TestTrackAdvanceProgressIgnoresAmbientTrack(t *testing.T) {
	track := Track{}
	track.advanceProgress(time.Now())
	if track.ProgressMs != nil || track.DurationMs != nil {
		t.Fatal("ambient timing must remain nil")
	}
}
```

- [ ] **Step 2: Run and verify failure**

Run from `apps/api-gql`:

```bash
go test ./internal/delivery/gql/now-playing-fetcher -run TestTrackAdvanceProgress -v
```

Expected: compilation fails because the timing fields and `advanceProgress` do not exist.

- [ ] **Step 3: Implement nullable timing and cache adjustment**

Extend `Track`:

```go
type Track struct {
	Artist             string    `json:"artist"`
	Title              string    `json:"title"`
	ImageUrl           string    `json:"image_url,omitempty"`
	ProgressMs         *int      `json:"progress_ms,omitempty"`
	DurationMs         *int      `json:"duration_ms,omitempty"`
	ProgressObservedAt time.Time `json:"progress_observed_at,omitempty"`
	fromCache          bool
}

func (t *Track) advanceProgress(now time.Time) {
	if t.ProgressMs == nil || t.DurationMs == nil || *t.DurationMs <= 0 || t.ProgressObservedAt.IsZero() {
		return
	}

	progress := *t.ProgressMs + int(now.Sub(t.ProgressObservedAt)/time.Millisecond)
	if progress < 0 {
		progress = 0
	}
	if progress > *t.DurationMs {
		progress = *t.DurationMs
	}
	t.ProgressMs = &progress
	t.ProgressObservedAt = now
}
```

After a successful Redis scan, call `cachedTrack.advanceProgress(time.Now())` before returning. When mapping a playing Spotify response, copy `ProgressMs` and `DurationMs` into local variables, store pointers, and set `ProgressObservedAt: time.Now()`. Leave Last.fm and VK timing nil.

- [ ] **Step 4: Run focused and package tests**

Run from `apps/api-gql`:

```bash
gofmt -w internal/delivery/gql/now-playing-fetcher/now-playing-fetcher.go internal/delivery/gql/now-playing-fetcher/now-playing-fetcher_test.go
```

Expected: all tests pass.

- [ ] **Step 5: Commit checkpoint, only with authorization**

```bash
git add apps/api-gql/internal/delivery/gql/now-playing-fetcher
git commit -m "feat(api): compensate cached track progress"
```

## Task 3: Extend Preset And GraphQL Contracts

**Files:**

- Modify: `libs/types/types/api/overlays/nowplaying.go`
- Modify: `apps/api-gql/internal/delivery/gql/schema/overlays/overlays.graphql`
- Create: `apps/api-gql/internal/delivery/gql/resolvers/now_playing_track.go`
- Create: `apps/api-gql/internal/delivery/gql/resolvers/now_playing_track_test.go`
- Modify: `apps/api-gql/internal/delivery/gql/resolvers/overlays.resolver.service.go`
- Regenerate: `libs/types/src/api.ts`
- Regenerate: `apps/api-gql/internal/delivery/gql/graph/generated.go`
- Regenerate: `apps/api-gql/internal/delivery/gql/gqlmodel/gqlmodel.go`

- [ ] **Step 1: Add a failing mapper test**

```go
package resolvers

import (
	"testing"

	nowplayingfetcher "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/now-playing-fetcher"
)

func TestMapNowPlayingTrackTiming(t *testing.T) {
	progress, duration := 70000, 238000
	got := mapNowPlayingTrack(&nowplayingfetcher.Track{
		Artist: "Glass Animals", Title: "Heat Waves", ImageUrl: "cover",
		ProgressMs: &progress, DurationMs: &duration,
	})
	if got.ProgressMs == nil || *got.ProgressMs != progress || got.DurationMs == nil || *got.DurationMs != duration {
		t.Fatalf("timing was not mapped: %#v", got)
	}

	ambient := mapNowPlayingTrack(&nowplayingfetcher.Track{Artist: "Artist", Title: "Title"})
	if ambient.ProgressMs != nil || ambient.DurationMs != nil {
		t.Fatalf("ambient timing must remain nil: %#v", ambient)
	}
}
```

- [ ] **Step 2: Add schema and enum source values**

Extend `NowPlayingOverlayPreset` and `NowPlayingOverlayTrack`:

```graphql
enum NowPlayingOverlayPreset {
	TRANSPARENT
	AIDEN_REDESIGN
	SIMPLE_LINE
	PULSE_STRIP
	AURA_STACK
	VINYL_HAZE
	SIGNAL_DECK
}

type NowPlayingOverlayTrack {
	artist: String!
	title: String!
	imageUrl: String
	progressMs: Int
	durationMs: Int
}
```

Add matching Go constants to `nowplaying.go`, include them in `AllPresets`, and return each value from `TSName()`:

```go
ChannelOverlayNowPlayingPresetPulseStrip ChannelOverlayNowPlayingPreset = "PULSE_STRIP"
ChannelOverlayNowPlayingPresetAuraStack  ChannelOverlayNowPlayingPreset = "AURA_STACK"
ChannelOverlayNowPlayingPresetVinylHaze  ChannelOverlayNowPlayingPreset = "VINYL_HAZE"
ChannelOverlayNowPlayingPresetSignalDeck ChannelOverlayNowPlayingPreset = "SIGNAL_DECK"
```

- [ ] **Step 3: Regenerate backend and shared TypeScript types**

Run from the repository root:

```bash
bun cli build gql
```

Run from `libs/types`:

```bash
bun run build
```

Expected: gqlgen models contain nullable `*int` timing fields and `src/api.ts` contains all seven preset values.

- [ ] **Step 4: Implement and use the mapper**

Create `now_playing_track.go`:

```go
package resolvers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	nowplayingfetcher "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/now-playing-fetcher"
)

func mapNowPlayingTrack(track *nowplayingfetcher.Track) *gqlmodel.NowPlayingOverlayTrack {
	if track == nil {
		return nil
	}

	var imageURL *string
	if track.ImageUrl != "" {
		imageURL = &track.ImageUrl
	}

	return &gqlmodel.NowPlayingOverlayTrack{
		Artist: track.Artist, Title: track.Title, ImageURL: imageURL,
		ProgressMs: track.ProgressMs, DurationMs: track.DurationMs,
	}
}
```

Replace the inline image/DTO block in `nowPlayingCurrentTrackSubscription` with `channel <- mapNowPlayingTrack(track)`.

- [ ] **Step 5: Run mapper and API tests**

Run from `apps/api-gql`:

```bash
gofmt -w internal/delivery/gql/resolvers/now_playing_track.go internal/delivery/gql/resolvers/now_playing_track_test.go internal/delivery/gql/resolvers/overlays.resolver.service.go
```

Expected: all focused tests pass.

- [ ] **Step 6: Commit checkpoint, only with authorization**

```bash
git add libs/types/types/api/overlays/nowplaying.go libs/types/src/api.ts apps/api-gql/internal/delivery/gql
git commit -m "feat(now-playing): add timing contract and presets"
```

## Task 4: Add Shared Progress State

**Files:**

- Modify: `libs/frontend-now-playing/src/types.ts`
- Create: `libs/frontend-now-playing/src/utils/progress.ts`
- Create: `libs/frontend-now-playing/tests/progress.test.ts`
- Create: `libs/frontend-now-playing/src/composables/use-track-progress.ts`

- [ ] **Step 1: Write failing pure timing tests**

```typescript
import { describe, expect, test } from 'bun:test'
import { formatTrackTime, normalizeTrackTiming } from '../src/utils/progress.ts'

describe('normalizeTrackTiming', () => {
	test('interpolates and clamps timed tracks', () => {
		expect(normalizeTrackTiming({ progressMs: 70_000, durationMs: 238_000 }, 5_000)).toEqual({
			mode: 'timed', progressMs: 75_000, durationMs: 238_000, percent: 75_000 / 238_000 * 100,
		})
		expect(normalizeTrackTiming({ progressMs: 237_000, durationMs: 238_000 }, 5_000).progressMs).toBe(238_000)
	})

	test('uses ambient mode for missing or invalid duration', () => {
		expect(normalizeTrackTiming({ progressMs: null, durationMs: null }, 0).mode).toBe('ambient')
		expect(normalizeTrackTiming({ progressMs: 0, durationMs: 0 }, 0).mode).toBe('ambient')
	})
})

test('formatTrackTime formats minute and hour durations', () => {
	expect(formatTrackTime(70_000)).toBe('01:10')
	expect(formatTrackTime(3_661_000)).toBe('1:01:01')
})
```

- [ ] **Step 2: Run and verify failure**

Run from the repository root:

```bash
bun test libs/frontend-now-playing/tests/progress.test.ts
```

Expected: module-not-found failure for `utils/progress.ts`.

- [ ] **Step 3: Add track fields and pure helpers**

Add nullable optional timing to `Track` and the four preset constants to `Preset`:

```typescript
export interface Track {
	artist: string
	title: string
	imageUrl?: string | null
	progressMs?: number | null
	durationMs?: number | null
}
```

Create `progress.ts` with exported `ProgressMode`, `NormalizedTrackTiming`, `normalizeTrackTiming(track, elapsedMs)`, and `formatTrackTime(ms)`:

```typescript
export type ProgressMode = 'timed' | 'ambient'
export interface TrackTimingInput { progressMs?: number | null; durationMs?: number | null }
export interface NormalizedTrackTiming {
	mode: ProgressMode
	progressMs: number
	durationMs: number
	percent: number
}

export function normalizeTrackTiming(track: TrackTimingInput, elapsedMs: number): NormalizedTrackTiming
export function formatTrackTime(ms: number): string
```

Clamp progress with `Math.min(Math.max(progress + elapsed, 0), duration)`. Return ambient mode with zero progress, duration, and percent when either value is absent or duration is non-positive.

- [ ] **Step 4: Run pure tests**

```bash
bun test libs/frontend-now-playing/tests/progress.test.ts
```

Expected: all timing tests pass.

- [ ] **Step 5: Add the lifecycle composable**

Create `use-track-progress.ts`:

```typescript
import { computed, onScopeDispose, ref, watch, type Ref } from 'vue'
import { formatTrackTime, normalizeTrackTiming } from '../utils/progress.js'
import type { Track } from '../types.js'

export function useTrackProgress(track: Ref<Track | null | undefined>) {
	const receivedAt = ref(Date.now())
	const clock = ref(Date.now())

	watch(track, () => {
		receivedAt.value = Date.now()
		clock.value = receivedAt.value
	}, { immediate: true })

	let interval: ReturnType<typeof setInterval> | undefined
	if (typeof window !== 'undefined') {
		interval = window.setInterval(() => { clock.value = Date.now() }, 1000)
	}
	onScopeDispose(() => {
		if (interval !== undefined) clearInterval(interval)
	})

	const timing = computed(() => normalizeTrackTiming(track.value ?? {}, clock.value - receivedAt.value))
	return {
		timing,
		elapsedLabel: computed(() => timing.value.mode === 'timed' ? formatTrackTime(timing.value.progressMs) : ''),
		durationLabel: computed(() => timing.value.mode === 'timed' ? formatTrackTime(timing.value.durationMs) : ''),
		progressStyle: computed(() => ({ '--track-progress': `${timing.value.percent}%` })),
	}
}
```

- [ ] **Step 6: Type-check the shared library source through a consumer build**

Run from `frontend/overlays`:

```bash
bun run build
```

Expected at this stage: existing app builds; new composable has no TypeScript errors.

- [ ] **Step 7: Commit checkpoint, only with authorization**

```bash
git add libs/frontend-now-playing/src/types.ts libs/frontend-now-playing/src/utils/progress.ts libs/frontend-now-playing/src/composables/use-track-progress.ts libs/frontend-now-playing/tests/progress.test.ts
git commit -m "feat(now-playing): add shared progress state"
```

## Task 5: Derive Artwork Palettes

**Files:**

- Create: `libs/frontend-now-playing/src/utils/palette.ts`
- Create: `libs/frontend-now-playing/tests/palette.test.ts`
- Create: `libs/frontend-now-playing/src/composables/use-artwork-palette.ts`

- [ ] **Step 1: Write failing palette tests**

```typescript
import { describe, expect, test } from 'bun:test'
import { contrastRatio, derivePalette, selectArtworkColors } from '../src/utils/palette.ts'

describe('artwork palette', () => {
	test('selects a frequent dominant and saturated accent', () => {
		const samples = [
			...Array.from({ length: 8 }, () => ({ r: 40, g: 70, b: 90 })),
			...Array.from({ length: 3 }, () => ({ r: 230, g: 40, b: 130 })),
		]
		const selected = selectArtworkColors(samples)
		expect(selected.dominant.b).toBeGreaterThan(selected.dominant.r)
		expect(selected.accent.r).toBeGreaterThan(selected.accent.g)
	})

	test('derives readable text and neutral fallback', () => {
		const palette = derivePalette([], 'rgba(0, 0, 0, 0)')
		expect(contrastRatio(palette.text, palette.surface)).toBeGreaterThanOrEqual(4.5)
		expect(palette.surface).toMatch(/^#[0-9a-f]{6}$/i)
	})
})
```

- [ ] **Step 2: Run and verify failure**

```bash
bun test libs/frontend-now-playing/tests/palette.test.ts
```

Expected: module-not-found failure for `utils/palette.ts`.

- [ ] **Step 3: Implement deterministic color helpers**

Create `palette.ts` with these public contracts:

```typescript
export interface Rgb { r: number; g: number; b: number }
export interface ArtworkPalette {
	surface: string
	surfaceAlt: string
	accent: string
	text: string
	mutedText: string
	glow: string
}

export const NEUTRAL_BASE: Rgb = { r: 32, g: 38, b: 42 }

export function selectArtworkColors(samples: Rgb[]): { dominant: Rgb; accent: Rgb }
export function derivePalette(samples: Rgb[], fallbackColor: string): ArtworkPalette
export function contrastRatio(foreground: string, background: string): number
```

Implementation rules:

- Quantize each channel to 32-value buckets and rank buckets by count for the dominant color.
- Rank accent candidates by HSL saturation multiplied by the square root of bucket count.
- Parse `#rgb`, `#rrggbb`, `rgb(...)`, and `rgba(...)`; treat zero alpha and invalid strings as `NEUTRAL_BASE`.
- Mix the dominant color toward `#101419` for `surface` and slightly less for `surfaceAlt`.
- Increase accent saturation only within valid RGB bounds.
- Choose and adjust near-white text until `contrastRatio(text, surface) >= 4.5`.
- Return lowercase six-digit hex values; return glow as `rgba(r, g, b, 0.35)`.

- [ ] **Step 4: Run palette tests**

```bash
bun test libs/frontend-now-playing/tests/palette.test.ts
```

Expected: all palette tests pass.

- [ ] **Step 5: Implement image loading and stale-result protection**

Create `use-artwork-palette.ts` with this public signature:

```typescript
export function useArtworkPalette(
	imageUrl: Ref<string | null | undefined>,
	fallbackColor: Ref<string>,
): { palette: Readonly<Ref<ArtworkPalette>>; paletteStyle: ComputedRef<Record<string, string>> }
```

Watch the image URL and fallback color. Increment a request token before every load, set `image.crossOrigin = 'anonymous'`, draw into a 32-by-32 canvas, skip pixels with alpha below 128, and call `derivePalette(samples, fallback)`. Catch load or canvas errors and call `derivePalette([], fallback)`. Before assigning the palette, verify the local token still matches the latest token.

Return both `palette` and this computed root style:

```typescript
const paletteStyle = computed(() => ({
	'--np-surface': palette.value.surface,
	'--np-surface-alt': palette.value.surfaceAlt,
	'--np-accent': palette.value.accent,
	'--np-text': palette.value.text,
	'--np-muted': palette.value.mutedText,
	'--np-glow': palette.value.glow,
}))
```

Guard browser-only APIs with `if (typeof window === 'undefined')` so the Nuxt dashboard can render without SSR errors.

- [ ] **Step 6: Run all shared helper tests**

```bash
bun test libs/frontend-now-playing/tests
```

Expected: timing and palette suites pass.

- [ ] **Step 7: Commit checkpoint, only with authorization**

```bash
git add libs/frontend-now-playing/src/utils/palette.ts libs/frontend-now-playing/src/composables/use-artwork-palette.ts libs/frontend-now-playing/tests/palette.test.ts
git commit -m "feat(now-playing): derive artwork palettes"
```

## Task 6: Build Pulse Strip And Aura Stack

**Files:**

- Create: `libs/frontend-now-playing/src/presets/pulse-strip.vue`
- Create: `libs/frontend-now-playing/src/presets/aura-stack.vue`
- Modify: `libs/frontend-now-playing/src/assets/style.css`

- [ ] **Step 1: Add shared visual primitives**

Extend `style.css` with:

```css
.np-progress {
	position: relative;
	overflow: hidden;
	height: clamp(6px, 1.2cqi, 10px);
	border-radius: 999px;
	background: color-mix(in srgb, var(--np-accent) 35%, var(--np-surface));
}

.np-progress::after {
	content: '';
	display: block;
	width: var(--track-progress, 0%);
	height: 100%;
	border-radius: inherit;
	background: var(--np-accent);
	transition: width 1s linear;
}

.np-progress[data-mode='ambient']::after {
	width: 42%;
	animation: np-ambient-progress 2.4s ease-in-out infinite;
}

@keyframes np-ambient-progress {
	from { transform: translateX(-120%); }
	to { transform: translateX(340%); }
}

@media (prefers-reduced-motion: reduce) {
	.np-progress::after { transition: none; }
	.np-progress[data-mode='ambient']::after { transform: none; animation: none; }
	.marque { animation: none; }
}
```

- [ ] **Step 2: Implement Pulse Strip**

Use `toRef(props, 'track')` with both shared composables. Measure a clipping container and an intrinsic-width metadata span with `useElementSize`; enable the existing `marque` class only when intrinsic width exceeds available width:

```typescript
const props = defineProps<{ track?: Track | null; settings: Settings }>()
const trackRef = toRef(props, 'track')
const metaViewport = ref<HTMLElement>()
const metaContent = ref<HTMLElement>()
const { width: viewportWidth } = useElementSize(metaViewport)
const { width: contentWidth } = useElementSize(metaContent)
const shouldMarquee = computed(() => contentWidth.value > viewportWidth.value)
const { timing, progressStyle } = useTrackProgress(trackRef)
const { paletteStyle } = useArtworkPalette(
	computed(() => props.track?.imageUrl),
	computed(() => props.settings.backgroundColor),
)
```

The root receives `paletteStyle`, metadata uses the existing font settings, and progress binds `timing.mode` plus `progressStyle`:

```vue
<template>
	<div v-if="track" class="pulse-strip" :style="paletteStyle">
		<img v-if="settings.showImage" class="pulse-strip__cover" :src="track.imageUrl ?? '/overlays/images/play.png'">
		<div ref="metaViewport" class="pulse-strip__meta-viewport">
			<div ref="metaContent" class="pulse-strip__meta" :class="{ marque: shouldMarquee }">
				<strong class="name">{{ track.title }}</strong><span aria-hidden="true">•</span><span class="artist">{{ track.artist }}</span>
			</div>
		</div>
		<div class="np-progress" :data-mode="timing.mode" :style="progressStyle" />
	</div>
</template>
```

Style the root as a three-column grid at wide widths and move progress to a full second row below 420px. Use `container-type: inline-size`, `minmax(0, 1fr)`, palette variables, and no fixed viewport units. Keep the viewport overflow hidden and the intrinsic metadata width at `max-content` so the measurement is meaningful.

- [ ] **Step 3: Implement Aura Stack**

Use one root grid with optional artwork, metadata card, and timing card:

```vue
<template>
	<div v-if="track" class="aura-stack" :data-mode="timing.mode" :style="paletteStyle">
		<img v-if="settings.showImage" class="aura-stack__cover" :src="track.imageUrl ?? '/overlays/images/play.png'">
		<section class="aura-stack__card">
			<h2 class="name">{{ track.title }}</h2>
			<p class="artist">{{ track.artist }}</p>
		</section>
		<section class="aura-stack__card aura-stack__timing">
			<span v-if="timing.mode === 'timed'">{{ elapsedLabel }}</span>
			<div class="np-progress" :data-mode="timing.mode" :style="progressStyle" />
			<span v-if="timing.mode === 'timed'">{{ durationLabel }}</span>
		</section>
	</div>
</template>
```

Use a maximum inline size near 360px, square cover, rounded surfaces, glow from `--np-glow`, and `.aura-stack[data-mode='ambient'] .aura-stack__timing { grid-template-columns: 1fr; }` so hidden labels leave no blank columns. Add `transition: background-color 300ms ease, color 300ms ease, box-shadow 300ms ease` to palette-colored surfaces.

- [ ] **Step 4: Build the overlays consumer**

Run from `frontend/overlays`:

```bash
bun run build
```

Expected: Vue and CSS compile without errors. These components are not selected until Task 8.

- [ ] **Step 5: Commit checkpoint, only with authorization**

```bash
git add libs/frontend-now-playing/src/assets/style.css libs/frontend-now-playing/src/presets/pulse-strip.vue libs/frontend-now-playing/src/presets/aura-stack.vue
git commit -m "feat(now-playing): add strip and stack presets"
```

## Task 7: Build Vinyl Haze And Signal Deck

**Files:**

- Create: `libs/frontend-now-playing/src/presets/vinyl-haze.vue`
- Create: `libs/frontend-now-playing/src/presets/signal-deck.vue`

- [ ] **Step 1: Implement Vinyl Haze**

Use the Aura Stack structure but derive muted component-local CSS variables with `color-mix()`:

```vue
<div
	v-if="track"
	class="vinyl-haze"
	:style="{
		...paletteStyle,
		'--haze-surface': 'color-mix(in srgb, var(--np-surface) 68%, #817536)',
		'--haze-accent': 'color-mix(in srgb, var(--np-accent) 52%, #d5cc55)',
	}"
>
```

Render optional square artwork, metadata card, and timing card. Apply `filter: saturate(.55) contrast(.92)` only to artwork, never to text. Use `--haze-surface` and `--haze-accent` for the cards and progress while retaining `--np-text` contrast.

- [ ] **Step 2: Implement Signal Deck**

Create a constant array of bar heights in script and render it with stable keys:

```typescript
const waveform = [10, 18, 14, 27, 20, 32, 17, 29, 24, 30, 19, 23, 16, 26, 18, 21, 13, 10]
```

Use this structure:

```vue
<template>
	<div v-if="track" class="signal-deck" :class="{ 'signal-deck--no-cover': !settings.showImage }" :style="paletteStyle">
		<img v-if="settings.showImage" class="signal-deck__cover" :src="track.imageUrl ?? '/overlays/images/play.png'">
		<section class="signal-deck__card signal-deck__meta">
			<h2 class="name">{{ track.title }}</h2><p class="artist">{{ track.artist }}</p>
		</section>
		<section class="signal-deck__card signal-deck__signal">
			<span v-if="timing.mode === 'timed'">{{ elapsedLabel }}</span>
			<div class="signal-deck__wave" aria-hidden="true">
				<i v-for="(height, index) in waveform" :key="index" :style="{ '--bar-height': `${height}px`, '--bar-index': index }" />
			</div>
			<span v-if="timing.mode === 'timed'">{{ durationLabel }}</span>
		</section>
		<div class="np-progress signal-deck__progress" :data-mode="timing.mode" :style="progressStyle" />
	</div>
</template>
```

Animate bars with scaleY and per-index negative delay. Disable bar animation under reduced motion. Use a two-column grid with artwork spanning both cards and progress spanning all columns. Below 480px, reduce artwork to 92px and hide every bar after index 11 with a container query. Set `.signal-deck--no-cover` to one column so both cards and progress span the full grid.

- [ ] **Step 3: Build the overlays consumer**

Run from `frontend/overlays`:

```bash
bun run build
```

Expected: all four new SFCs compile.

- [ ] **Step 4: Commit checkpoint, only with authorization**

```bash
git add libs/frontend-now-playing/src/presets/vinyl-haze.vue libs/frontend-now-playing/src/presets/signal-deck.vue
git commit -m "feat(now-playing): add haze and signal presets"
```

## Task 8: Wire Renderer, Dashboard, And OBS Consumers

**Files:**

- Modify: `libs/frontend-now-playing/src/now-playing.vue`
- Modify: `frontend/overlays/src/composables/now-playing/use-now-playing-socket.ts`
- Modify: `web/layers/dashboard/pages/dashboard/overlays/now-playing.vue`
- Modify: `web/layers/dashboard/pages/dashboard/overlays/now-playing/now-playing-form.vue`
- Regenerate: `frontend/overlays/src/gql/*`
- Regenerate: `web/app/gql/*`

- [ ] **Step 1: Select the new preset components**

Import all four `.vue` files with extensions and extend the existing switch:

```typescript
case Preset.PULSE_STRIP:
	return PresetPulseStrip
case Preset.AURA_STACK:
	return PresetAuraStack
case Preset.VINYL_HAZE:
	return PresetVinylHaze
case Preset.SIGNAL_DECK:
	return PresetSignalDeck
```

Keep Transparent as the default.

- [ ] **Step 2: Request timing in both subscriptions**

Add these fields after `imageUrl` in the overlays app and dashboard subscription documents:

```graphql
progressMs
durationMs
```

In the dashboard computed mapping, pass both values through. Extend the sample track:

```typescript
const defaultNowPlayingTrack = {
	imageUrl: 'https://i.scdn.co/image/ab67616d0000b273e7fbc0883149094912559f2c',
	artist: 'Slipknot',
	title: 'Psychosocial',
	progressMs: 70_000,
	durationMs: 238_000,
}
```

- [ ] **Step 3: Add dashboard select labels**

Add these `SelectItem` entries after the existing presets:

```vue
<SelectItem value="PULSE_STRIP">Pulse Strip</SelectItem>
<SelectItem value="AURA_STACK">Aura Stack</SelectItem>
<SelectItem value="VINYL_HAZE">Vinyl Haze</SelectItem>
<SelectItem value="SIGNAL_DECK">Signal Deck</SelectItem>
```

- [ ] **Step 4: Regenerate frontend GraphQL clients**

Run from `frontend/overlays`:

```bash
bun run codegen
```

Run from `web`:

```bash
bun run graphql-codegen
```

Expected: generated operation result types expose nullable progress and duration and the preset enum includes all values.

- [ ] **Step 5: Run helper tests and both frontend builds**

Run from the repository root:

```bash
bun test libs/frontend-now-playing/tests
```

Run from `frontend/overlays`:

```bash
bun run build
```

Run from `web`:

```bash
bun run build
```

Expected: tests pass and both production builds finish without type or GraphQL errors.

- [ ] **Step 6: Commit checkpoint, only with authorization**

```bash
git add libs/frontend-now-playing frontend/overlays/src/composables/now-playing frontend/overlays/src/gql web/layers/dashboard/pages/dashboard/overlays/now-playing.vue web/layers/dashboard/pages/dashboard/overlays/now-playing web/app/gql
git commit -m "feat(now-playing): add adaptive artwork presets"
```

## Task 9: End-To-End Verification

**Files:**

- Verify only; fix failures in the files owned by Tasks 1 through 8.

- [ ] **Step 1: Run backend tests**

Run from `libs/integrations`:

```bash
go test ./spotify -v
```

Run from `apps/api-gql`:

```bash
go test ./internal/delivery/gql/now-playing-fetcher ./internal/delivery/gql/resolvers -v
```

Expected: all packages pass.

- [ ] **Step 2: Run frontend tests, builds, and lint**

Run from the repository root:

```bash
bun test libs/frontend-now-playing/tests
bun lint
```

Run `bun run build` in `frontend/overlays` and `web`.

Expected: tests, lint, and builds pass.

- [ ] **Step 3: Verify dashboard preview visually**

Open the Caddy-served dashboard URL from `SITE_BASE_URL` or `http://localhost:3005`, navigate to the now-playing overlay settings, and verify:

- All four labels select the intended component.
- Palette follows the sample or current artwork and text remains readable.
- Pulse Strip and both vertical cards show real sample progress.
- Signal Deck shows elapsed time, animated waveform, total time, and separate progress.
- `showImage=false` closes the artwork gap.
- Long title and artist remain contained.
- Narrow and wide preview containers preserve hierarchy.

- [ ] **Step 4: Verify source modes on the public overlay route**

With Spotify enabled, confirm subscription corrections never move progress backward, elapsed time advances locally, and Signal Deck shows waveform plus progress. With Last.fm or VK data, confirm timing labels are absent and progress uses ambient sweep. Test a missing artwork URL and transparent `backgroundColor` to confirm the neutral fallback.

- [ ] **Step 5: Verify reduced motion**

Emulate `prefers-reduced-motion: reduce`. Confirm waveform and ambient sweep stop, determinate timing remains correct, and metadata remains visible.

- [ ] **Step 6: Inspect the final diff**

```bash
git status --short
```

Expected: only intended source, tests, generated GraphQL/type files, the design spec, and this plan are changed; `git diff --check` prints no errors.

- [ ] **Step 7: Final commit checkpoint, only with authorization**

If authorized and verification fixes remain uncommitted:

```bash
git add <only-the-verified-feature-files>
git commit -m "fix(now-playing): polish adaptive presets"
```
