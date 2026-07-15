# Song Request Overlay Styles Design

## Goal

Add persisted, realtime-configurable visual styles for the public song-request overlay while preserving the existing YouTube playback lifecycle, transparent paused state, and dashboard controls.

## Scope

The feature adds six overlay styles:

- `CINEMA`: current full-size video with a bottom status bar.
- `COMPACT`: small vertical video card with metadata below.
- `TICKER`: bottom line with configurable colors and readable marquee text.
- `STUDIO`: horizontal cover card with metadata, time, decorative waveform, and progress.
- `PORTRAIT`: vertical cover card with metadata and progress.
- `PILL`: minimal single-row cover, title/requester, and progress.

It also adds a dedicated dashboard modal, local previews, persisted settings, realtime updates, migration of `hide_on_pause`, and missing English/Russian translations.

## Data Model

Create a dedicated `channels_song_requests_overlay_settings` table through the project migration CLI.

Columns:

| Column | Type | Rules | Default |
| --- | --- | --- | --- |
| `id` | UUID | primary key | `uuidv7()` |
| `channel_id` | TEXT | unique, not null, FK to `channels(id)`, cascade delete | none |
| `style` | `song_request_overlay_style` | not null | `CINEMA` |
| `accent_color` | VARCHAR(9) | `#RRGGBB` or `#RRGGBBAA` | `#8B5CF6` |
| `ticker_background_color` | VARCHAR(9) | `#RRGGBB` or `#RRGGBBAA` | `#111827E6` |
| `ticker_text_color` | VARCHAR(9) | `#RRGGBB` or `#RRGGBBAA` | `#FFFFFF` |
| `ticker_speed` | INT | 10 through 100 px/s | 35 |
| `hide_on_pause` | BOOLEAN | not null | true |
| `created_at` | TIMESTAMPTZ | not null | `now()` |
| `updated_at` | TIMESTAMPTZ | not null | `now()` |

The migration copies each existing `channels_song_requests_settings.hide_on_pause` value into the new table, then removes the old column. The down migration restores the old column and values before removing the new table and enum.

`player_no_cookie_mode` remains in `channels_song_requests_settings`; this change does not redesign YouTube privacy mode.

## Backend Architecture

New code follows the project entity and pgx repository pattern:

- Entity: typed overlay settings and style enum with defaults.
- Repository interface: get by channel ID and upsert.
- pgx implementation: maps DB rows to the entity.
- API service: default fallback, validation, persistence, and realtime publication.
- GraphQL mapper: entity to GraphQL DTO.

No new GORM model is introduced.

The service returns defaults without inserting a row when settings are absent. The first dashboard Save performs an upsert.

## GraphQL Contract

Add `SongRequestOverlayStyle` values:

- `CINEMA`
- `COMPACT`
- `TICKER`
- `STUDIO`
- `PORTRAIT`
- `PILL`

Add `SongRequestOverlaySettings` and matching update input containing style, accent color, ticker colors, ticker speed, and hide-on-pause.

Operations:

- Authenticated dashboard query using the selected dashboard channel.
- Authenticated update mutation requiring `MANAGE_SONG_REQUESTS`.
- Public `@noRateLimit` subscription accepting `apiKey`.

The public subscription validates the API key, resolves the channel, sends current/default settings immediately, then forwards updates from `WsRouter`.

## Realtime Flow

The service uses a channel-scoped key such as `songrequests:overlay-settings:<channelId>`.

Save flow:

1. Validate enum, colors, speed, and channel permission.
2. Upsert settings through the pgx repository.
3. Invalidate any settings cache used by the service.
4. Publish the persisted DTO through `WsRouter`.
5. Return success to the dashboard.

The database is the source of truth. A publication failure does not roll back a successful database write; reconnecting the overlay receives the persisted initial state.

## Dashboard UX

Add a settings icon beside the Overlay link controls in the song-request links card. It opens a dedicated overlay settings modal, separate from the general song-request settings modal.

The modal contains:

- Six selectable style preview cards.
- A larger live local preview using sample title/requester data.
- Shared accent-color input.
- Ticker-only background color, text color, and scroll-speed inputs.
- Hide-on-pause switch.
- Cancel and Save buttons.

Preview changes are local. The real overlay changes only after Save. Save shows a success toast, updates query cache, and closes the modal.

The old hide-on-pause field is removed from the general song-request modal and GraphQL settings input/output.

## Shared Renderer Architecture

Place reusable renderer components under `web/app/components/song-request-overlay/` so both the public overlay and dashboard preview use the same visuals.

Components:

- A renderer that selects a style from the enum.
- One focused component per style.
- Shared typed props for title, requester, video ID, position, duration, playing state, and colors.
- A media slot for styles that display video or preview imagery.

The public overlay owns the YouTube API player. The preview never creates a YouTube player and uses the YouTube thumbnail instead.

## Style Behavior

### Cinema

- Fills the browser-source viewport.
- Shows the actual YouTube video.
- Keeps only the bottom status bar opaque.
- Displays title, requester, time, and progress.

### Compact

- Bottom-left card up to 480 px wide.
- Actual 16:9 video above metadata.
- Metadata contains title, requester, time, and progress.

### Ticker

- Bottom full-width line.
- Text format: `@requester • title`.
- Uses configurable background and text colors.
- Uses the shared accent for separators/progress detail.
- Text remains static when it fits.
- Overflowing text uses a seamless right-to-left marquee.
- Animation duration derives from measured travel distance divided by configured pixels per second.

### Studio

- Bottom-left horizontal card up to 680 px wide.
- Square YouTube thumbnail on the left.
- Title and `Requested by ...` on the right.
- Bottom row contains elapsed time, decorative waveform, total time, and progress.
- Decorative waveform animates only during playback.

### Portrait

- Bottom-left vertical card up to 360 px wide.
- Large thumbnail, title, requester, and time/progress blocks.
- Uses the accent for progress and a restrained glow.

### Pill

- Bottom full-width or container-width row around 56 px tall.
- Miniature thumbnail.
- Text format: `title • @requester`.
- Compact progress appears on the right.

## YouTube Player Lifecycle

Cinema and Compact display the actual iframe. Other styles display a thumbnail or no image while keeping the same iframe mounted as an invisible audio host.

The iframe must never be removed or changed to `display: none` during pause/style changes. Visibility uses opacity and pointer-events so the player retains dimensions and pause/resume remains reliable.

When `hide_on_pause` is true, the visual renderer becomes transparent while the player stays mounted. When false, the selected style remains visible in a paused state.

Changing style or color must not recreate the YouTube player or reset playback position.

## Validation and Error Handling

- Backend and database validate speed range.
- Backend validates colors as `#RRGGBB` or `#RRGGBBAA`.
- Unknown or missing style falls back to Cinema.
- Missing settings fall back to documented defaults.
- Overlay subscription failure leaves playback functional with defaults.
- Dashboard mutation errors keep the modal open and show an error toast.
- Ticker measurement reacts to viewport and text-size changes.

## Localization

Add English and Russian translations for:

- Existing `songRequests.links.*` labels currently rendering as raw keys.
- Overlay settings button, modal title, labels, descriptions, Save/Cancel feedback.
- Six style names and descriptions.
- Preview labels and Ticker controls.

Other locales may use the existing English fallback behavior.

## Verification

No automated test files are added, per the current task constraint.

Verification includes:

- Migration Up/Down review.
- GraphQL code generation.
- Go build for affected workspace modules.
- Frontend GraphQL generation and available type checking.
- Playwright checks for all six modes.
- Playwright Save-to-realtime-update check with an already-open overlay.
- Pause/resume and style-switch checks without iframe recreation.
- Transparent background and long Ticker text checks.
- Reconnect check confirming persisted initial settings.
