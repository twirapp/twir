# Now Playing Overlay Presets Design

## Goal

Add four adaptive now-playing overlay presets based on the supplied references. Their colors follow the current artwork, Spotify tracks display real playback progress, and sources without timing data use an animated ambient presentation.

The existing presets and now-playing behavior remain available.

## Scope

Add these presets:

- `PULSE_STRIP`: compact horizontal artwork, metadata, and progress.
- `AURA_STACK`: large vertical artwork with separate metadata and timing cards.
- `VINYL_HAZE`: a muted, vintage variation of the vertical card composition.
- `SIGNAL_DECK`: horizontal artwork, metadata card, animated waveform, timing, and progress.

The feature also adds Spotify timing to the existing now-playing subscription, artwork-derived palettes for the four new presets, an ambient fallback for Last.fm and VK, responsive layouts, dashboard selection, and matching dashboard preview behavior.

This change does not add user-configurable palette fields, audio analysis, or a database migration.

## Preset Behavior

### Pulse Strip

- Render a compact single-row card.
- Place a small square artwork image first, followed by `title • artist` and a progress bar.
- Truncate safely at very narrow widths and use marquee only when the measured text overflows.
- When artwork is hidden, let metadata use the released space.

### Aura Stack

- Render a large square artwork image above two rounded cards.
- Use the first card for title and artist.
- Use the second card for elapsed time, progress, and total time.
- Apply a restrained glow derived from the artwork accent.
- When artwork is hidden, retain the two information cards without an empty image area.

### Vinyl Haze

- Use the same structural hierarchy as Aura Stack.
- Desaturate and soften the derived palette to produce the muted reference treatment.
- Apply a subtle haze to the artwork and background without reducing text contrast.
- Keep timing and ambient behavior identical to Aura Stack.

### Signal Deck

- Render a horizontal layout with square artwork on the left and two cards on the right.
- Use the upper card for title and artist.
- Use the lower card for elapsed time, an animated decorative waveform, and total time.
- Render a separate full-width progress bar below both columns.
- Keep the waveform active for Spotify as well as ambient sources. It is decorative and does not claim to represent audio amplitude.
- On narrow containers, reduce artwork size and waveform density while preserving the same hierarchy.

## Timed And Ambient Modes

The shared renderer derives one of two modes from the track payload.

### Timed Mode

Timed mode is active when both `progressMs` and a positive `durationMs` are present.

- Display elapsed and total time where the preset includes timing labels.
- Clamp progress to the inclusive range from zero through duration.
- Advance elapsed time locally between subscription updates.
- Correct the local baseline whenever a new subscription payload arrives.
- Stop local advancement at duration.
- Signal Deck also runs its decorative waveform.

### Ambient Mode

Ambient mode is active when either timing value is absent or duration is not positive. Last.fm and VK intentionally use this mode.

- Hide numeric timing labels rather than displaying invented values.
- Replace determinate progress with a soft repeating sweep.
- Continue the Signal Deck waveform animation.
- Under `prefers-reduced-motion`, replace the sweep and waveform motion with stable decorative states.

## Dynamic Artwork Palette

Palette extraction runs in `@twir/frontend-now-playing` and only affects the four new presets. Existing presets retain their current color behavior.

For each new artwork URL:

1. Load the image with anonymous CORS enabled.
2. Draw it into a small canvas to bound CPU and memory usage.
3. Ignore transparent pixels and quantize the remaining samples into color buckets.
4. Select a representative dominant color and a saturated accent candidate.
5. Derive dark surface, secondary surface, accent, primary text, and muted text colors.
6. Adjust text colors to maintain readable contrast against their surfaces.
7. Expose the result as CSS custom properties on the preset root.

Palette changes transition softly when the track changes. A monotonically increasing request token prevents a slow previous image load from replacing the current track's palette.

If artwork is missing, cannot be loaded, or cannot be read because of CORS, derive the palette from the existing `backgroundColor` setting. A transparent or invalid fallback color uses a documented neutral dark base. The artwork may still drive the palette when `showImage` is false because visual visibility and color extraction are independent settings.

## Spotify Timing Data

Extend both Spotify response paths in `libs/integrations/spotify/track.go`:

- `/v1/me/player/currently-playing`
- `/v1/me/player`

Both responses map root `progress_ms` and item `duration_ms` into `GetTrackResponse`.

The now-playing fetcher's internal `Track` adds nullable progress and duration plus the time at which progress was observed. Spotify sets these fields; Last.fm and VK leave them nil.

The current Redis cache stores the observation time with the Spotify track. When returning a cached timed track, the fetcher adds elapsed wall-clock time to the observed position and clamps it to duration. This prevents the ten-second cache from repeatedly returning an old position. Non-timed tracks are unchanged.

## GraphQL Contract

Add nullable fields to `NowPlayingOverlayTrack`:

```graphql
progressMs: Int
durationMs: Int
```

The subscription resolver maps Spotify timing after cache compensation. It emits null timing for Last.fm and VK.

Add these values to every now-playing preset declaration and generated representation:

- `PULSE_STRIP`
- `AURA_STACK`
- `VINYL_HAZE`
- `SIGNAL_DECK`

The persisted `preset` column is already string-compatible, so no migration is required. GraphQL code must be regenerated with `bun cli build gql` after the schema edit.

## Shared Renderer Architecture

Keep presentation in `libs/frontend-now-playing` so the OBS overlay and dashboard preview use identical components.

Add:

- One Vue component for each new preset.
- A palette composable that owns image loading, sampling, fallback, and CSS color values.
- A progress composable that owns mode selection, clamping, local interpolation, formatting, and reduced-motion behavior.
- Small pure color and timing helpers where needed for deterministic tests.

Update the existing renderer switch to select all seven presets. Unknown values continue to fall back to the existing Transparent preset.

The OBS overlay and dashboard preview subscriptions request `progressMs` and `durationMs`. The dashboard's sample track includes stable timing values so all timed UI can be previewed when no integration is connected.

## Responsive Layout

Each preset owns a transparent root container and adapts to its available browser-source width rather than assuming the reference image dimensions.

- Use fluid dimensions with `clamp()` for artwork, spacing, and typography.
- Use container queries for structural adjustments where width alone is insufficient.
- Keep vertical presets bounded to a readable card width while allowing proportional growth.
- Collapse or reduce horizontal artwork and waveform density at narrow widths.
- Preserve text ellipsis or measured marquee behavior for long metadata.
- Avoid fixed viewport assumptions so dashboard previews and OBS sources render the same component correctly.

Existing `hideTimeout` behavior remains outside the preset components and continues to hide the rendered track after the configured interval.

## Animation And Performance

- Update displayed elapsed text once per second.
- Drive determinate progress from the locally interpolated position and use a linear CSS transition between updates.
- Build the Signal Deck waveform from CSS bars with varied heights, phases, and durations.
- Use CSS-only ambient sweeps after mode selection; do not issue extra network requests.
- Pause component timers and discard pending palette results when the component unmounts or track data becomes unavailable.
- Keep palette sampling resolution small and perform it only when the artwork URL or fallback background changes.

## Dashboard UX

Add the four new labels to the existing Style select:

- Pulse Strip
- Aura Stack
- Vinyl Haze
- Signal Deck

No additional form fields are introduced. `backgroundColor` remains the fallback palette source, `showImage` controls artwork visibility, and existing font family/weight settings apply to text in every new preset.

The dashboard preview uses the selected preset immediately through the existing local form state. Saving persists the enum through the existing mutation.

## Error Handling

- Treat missing or non-positive Spotify duration as ambient mode.
- Clamp negative progress to zero and progress beyond duration to duration.
- Ignore stale image load completions after a track change.
- Fall back to `backgroundColor` when image loading, CORS, canvas access, or sampling fails.
- Fall back to the neutral dark base when `backgroundColor` is transparent or invalid.
- Preserve the current Transparent fallback for unknown preset values.
- Do not let palette or animation failures hide otherwise valid track metadata.

## Verification

Automated checks cover:

- Mapping `progress_ms` and `duration_ms` from both Spotify API response shapes.
- Redis cache compensation and duration clamping.
- Nullable GraphQL mapping for non-Spotify tracks.
- Timed-versus-ambient mode selection, negative and overflow clamping, and time formatting.
- Palette derivation from sampled colors and fallback behavior for missing/invalid input.

Build and integration verification includes:

- `bun cli build gql` after schema changes.
- Focused Go tests for the Spotify integration and now-playing fetcher.
- Available frontend type checking and builds for the shared library, overlays app, and web app.
- Dashboard preview checks for all four new preset selections.
- OBS route checks at wide and narrow viewport sizes.
- Spotify checks confirming real elapsed/total timing, smooth correction after subscription updates, and Signal Deck waveform plus progress.
- Last.fm/VK checks confirming no numeric timing and the ambient animation state.
- Missing artwork, failed palette extraction, `showImage=false`, long metadata, transparent fallback color, and reduced-motion checks.

## Out Of Scope

- Real audio waveform or amplitude analysis.
- User-editable surface, text, accent, or waveform colors.
- Changes to Spotify polling frequency.
- Changes to integration priority or paused-track fallback behavior.
- Redesigning the three existing presets.
