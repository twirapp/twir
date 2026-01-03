# IMAGE Layer Type Feature

## Overview

Added support for IMAGE layer type to custom overlays, allowing users to display images from URLs in their stream overlays.

## Changes Made

### 1. Database

- **Migration**: `libs/migrations/postgres/20260103024719_add_image_layer_type.sql`
  - Added `IMAGE` value to `channels_overlays_layers_type` enum

### 2. Backend (Go)

#### GraphQL Schema

- **File**: `apps/api-gql/internal/delivery/gql/schema/overlays/overlays-custom.graphql`
  - Added `IMAGE` to `ChannelOverlayLayerType` enum
  - Added `imageUrl: String!` field to `ChannelOverlayLayerSettings` type
  - Added `imageUrl: String!` field to `ChannelOverlayLayerSettingsInput` input type

#### Models & Entities

- **File**: `libs/repositories/channels_overlays/model/model.go`
  - Added `OverlayTypeIMAGE` constant
  - Added `ImageUrl` field to `OverlayLayerSettings` struct

- **File**: `apps/api-gql/internal/entity/channel_overlay.go`
  - Added `ChannelOverlayTypeIMAGE` constant
  - Added `ImageUrl` field to `ChannelOverlayLayerSettings` struct

#### Mappers

- **File**: `apps/api-gql/internal/delivery/gql/mappers/channel_overlay.go`
  - Added IMAGE type mapping in `ChannelOverlayLayerTypeEntityToGql`
  - Added IMAGE type mapping in `ChannelOverlayLayerTypeGqlToEntity`
  - Added `ImageURL` field mapping in `ChannelOverlayLayerSettingsEntityToGql`

#### Services & Resolvers

- **File**: `apps/api-gql/internal/services/channels_overlays/channels_overlays.go`
  - Added `ImageUrl` field mapping in `modelToEntity` function
  - Added `ImageUrl` field mapping in `Update` function

- **File**: `apps/api-gql/internal/delivery/gql/resolvers/overlays-custom.resolver.go`
  - Added `ImageUrl` field mapping in `ChannelOverlayCreate` resolver
  - Added `ImageUrl` field mapping in `ChannelOverlayUpdate` resolver

### 3. Frontend - Dashboard

#### Types & Composables

- **File**: `frontend/dashboard/src/features/overlay-builder/types/index.ts`
  - Added `imageUrl?: string` to `LayerSettings` interface

- **File**: `frontend/dashboard/src/features/overlay-builder/composables/useOverlayBuilder.ts`
  - Added default settings for IMAGE layer type in `addLayer` function
  - Set `periodicallyRefetchData` to `false` for IMAGE layers (no need to poll)
  - Default placeholder image: `https://via.placeholder.com/300x200`

#### UI Components

- **File**: `frontend/dashboard/src/features/overlay-builder/OverlayBuilder.vue`
  - Added `addImageLayer` function
  - Added IMAGE layer button to "Add Layer" dialog with 🖼️ emoji

- **File**: `frontend/dashboard/src/features/overlay-builder/components/layer-editors/ImageLayerEditor.vue` (NEW)
  - Image URL input field
  - Live preview of the image
  - "Use Placeholder Image" button
  - Info tip about using Twir variables in URL

- **File**: `frontend/dashboard/src/features/overlay-builder/components/ImageLayerPreview.vue` (NEW)
  - Renders image on canvas
  - Shows error state if image fails to load
  - Shows placeholder state if no URL provided

- **File**: `frontend/dashboard/src/features/overlay-builder/components/PropertiesPanel.vue`
  - Added `ImageLayerEditor` component for IMAGE layer type

- **File**: `frontend/dashboard/src/features/overlay-builder/components/Canvas.vue`
  - Added `ImageLayerPreview` import and rendering for IMAGE layers

#### Data Handling

- **File**: `frontend/dashboard/src/components/registry/overlays/edit.vue`
  - Added `imageUrl` field when loading overlay data from GraphQL
  - Added `imageUrl` field when converting to mutation input

- **File**: `frontend/dashboard/src/components/registry/overlays/helpers.ts`
  - Added `IMAGE` case to `convertOverlayLayerTypeToText` helper (returns "Image")

### 4. Frontend - Overlays (Renderer)

#### Components

- **File**: `frontend/overlays/src/components/image-layer.vue` (NEW)
  - Renders IMAGE layer with proper positioning and rotation
  - Uses `object-fit: contain` to preserve aspect ratio
  - Absolute positioning with transform for rotation support

- **File**: `frontend/overlays/src/pages/overlays.vue`
  - Added `imageLayer` component import
  - Added conditional rendering for IMAGE layer type

#### Composables

- **File**: `frontend/overlays/src/composables/overlays/use-overlays.ts`
  - Added `imageUrl: string` to `LayerSettings` interface
  - Added `imageUrl` field mapping when transforming GraphQL data

## Usage

### Creating an IMAGE Layer

1. Open Overlay Editor in Dashboard
2. Click "Add Layer" button
3. Select "Image Layer" with 🖼️ icon
4. In Properties Panel, enter image URL in "Image URL" field
5. See live preview of the image
6. Adjust position, size, and rotation as needed
7. Click Save

### Image URL Features

- Direct URLs to images (PNG, JPG, GIF, WebP, etc.)
- Can use Twir variables: `$(user.login)`, `$(stream.title)`, etc.
- Example: `https://example.com/avatars/$(user.login).png`

### Properties

IMAGE layers support all standard layer properties:

- Position (X, Y)
- Size (Width, Height)
- Rotation (0-360°)
- Opacity (0-100%)
- Visibility toggle
- Lock/Unlock

## Technical Details

### Default Settings

When creating a new IMAGE layer:

- Width: 200px
- Height: 200px
- Default URL: `https://via.placeholder.com/300x200`
- `periodicallyRefetchData`: `false` (no polling needed for static images)

### Rendering

- **Dashboard**: Uses `<img>` tag with `object-fit: contain` in preview
- **Overlays**: Uses `<img>` tag with absolute positioning and CSS transforms for rotation
- Error handling: Shows error icon if image fails to load
- Fallback: Shows placeholder icon if no URL provided

## Migration

To apply the database changes:

```bash
bun cli m run --skip-clickhouse
```

To regenerate GraphQL types:

```bash
# Backend
bun cli build gql

# Frontend Dashboard
cd frontend/dashboard && bun run codegen

# Frontend Overlays
cd frontend/overlays && bun run codegen
```

## Testing Checklist

- [x] Database migration applies successfully
- [x] Backend compiles without errors
- [x] GraphQL schema updated
- [x] Can create IMAGE layer in dashboard
- [x] Can edit image URL
- [x] Image preview shows in editor
- [x] Can adjust position, size, rotation
- [x] Can save overlay with IMAGE layer
- [x] IMAGE layer renders correctly in overlay viewer
- [x] Rotation works correctly
- [x] Multiple IMAGE layers can be added
- [x] IMAGE layers work alongside HTML layers

## Bug Fixes

### Issue: Image URL not updating/saving

**Problem**: When editing the image URL in the dashboard, the preview didn't update, and the URL wasn't being saved to the database.

**Root Causes**:

1. **Missing GraphQL field**: The `imageUrl` field was not included in GraphQL queries, so it wasn't being fetched from or sent to the server.
2. **Missing initialization**: When loading existing overlays, the `imageUrl` field wasn't being copied from the server response into the builder state.

**Fixes Applied**:

1. **File**: `frontend/dashboard/src/api/overlays/custom.ts`
   - Added `imageUrl` field to all GraphQL queries (channelOverlaysQuery, channelOverlayByIdQuery, channelOverlayCreateMutation, channelOverlayUpdateMutation)

2. **File**: `frontend/dashboard/src/features/overlay-builder/OverlayBuilder.vue`
   - Added `imageUrl: layer.settings?.imageUrl || ''` to the `loadInitialProject` function to properly load imageUrl from server data

3. **File**: `frontend/overlays/src/composables/overlays/use-custom-overlay.ts`
   - Added `imageUrl` field to the overlay query for the renderer

4. **Regenerated GraphQL types** in both dashboard and overlays frontends

**Result**: Image URLs now properly save to the database, load from the database, and update in real-time when edited.

## Notes

- IMAGE layers do not support `periodicallyRefetchData` feature (not needed for images)
- HTML layer settings fields are still required in GraphQL schema but left empty for IMAGE layers
- Image loading is handled by browser's native `<img>` tag
- No CORS restrictions - images must be publicly accessible
- For dynamic images, Twir variable parsing happens server-side when the overlay loads

## Troubleshooting

If images are not showing or saving:

1. Check browser console for GraphQL errors
2. Verify that `imageUrl` field is present in GraphQL queries
3. Ensure GraphQL types are regenerated after schema changes
4. Confirm that the image URL is accessible (no CORS issues)
5. Check that settings object includes all required fields (htmlOverlayHtml, htmlOverlayCss, htmlOverlayJs, htmlOverlayDataPollSecondsInterval, imageUrl)
