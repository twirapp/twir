# IMAGE Layer Debugging Guide

## Issue Description

The IMAGE layer is being saved but displays "failed to load image" on the canvas in the dashboard.

## Changes Made

### 1. Fixed ImageLayerPreview Component Validation

**File**: `frontend/dashboard/src/features/overlay-builder/components/ImageLayerPreview.vue`

- Added `hasValidUrl` computed property to check if imageUrl is not empty/null/undefined
- Changed condition from `v-if="imageUrl && !imageError"` to `v-if="hasValidUrl && !imageError"`
- This prevents attempting to load empty strings as images

### 2. Fixed Overlay Renderer Validation

**File**: `frontend/overlays/src/components/image-layer.vue`

- Added `hasValidUrl` computed property similar to dashboard component
- Prevents rendering empty image URLs in the actual overlay

### 3. Fixed PropertiesPanel Emit

**File**: `frontend/dashboard/src/features/overlay-builder/components/PropertiesPanel.vue`

- Changed `@update="handleUpdateLayerProperties"` to `@update="emit('update', $event)"`
- Fixed missing function reference that was preventing updates

### 4. Improved Layer Settings Initialization

**File**: `frontend/dashboard/src/features/overlay-builder/composables/useOverlayBuilder.ts`

- Reordered settings for IMAGE layers to put `imageUrl` first (more logical order)
- Ensures new IMAGE layers always have a placeholder URL by default

### 5. Added Comprehensive Debug Logging

Added console logging to track imageUrl throughout the data flow:

#### `frontend/dashboard/src/features/overlay-builder/OverlayBuilder.vue`
- Logs each layer being loaded with its settings
- Shows imageUrl value for debugging

#### `frontend/dashboard/src/features/overlay-builder/components/Canvas.vue`
- Logs when IMAGE layer component is mounted
- Shows imageUrl being passed to ImageLayerPreview

#### `frontend/dashboard/src/features/overlay-builder/components/ImageLayerPreview.vue`
- Logs component mount with initial imageUrl
- Logs imageUrl changes
- Logs successful image loads
- Logs image load errors with the URL that failed

#### `frontend/dashboard/src/components/registry/overlays/edit.vue`
- Logs raw API data when loading overlay
- Logs layer conversion from API format
- Logs save operation with layer data
- Logs API format conversion before sending

## How to Debug

### Step 1: Open Browser DevTools Console

Open the browser console (F12) and filter for:
- `[edit.vue]` - API data loading/saving
- `[OverlayBuilder]` - Project data loading
- `[Canvas]` - Layer rendering
- `[ImageLayerPreview]` - Image component lifecycle

### Step 2: Create or Load an IMAGE Layer

1. Create a new IMAGE layer or load an existing overlay with IMAGE layers
2. Watch the console for log messages

Expected flow:
```
[edit.vue] Loading overlay data from API: {...}
[edit.vue] Converting layer from API: { type: 'IMAGE', imageUrl: '...', ... }
[OverlayBuilder] Loading layer 0: { type: 'IMAGE', imageUrl: '...', ... }
[Canvas] IMAGE layer mounted: layer-xxx imageUrl: '...'
[ImageLayerPreview] Component mounted with imageUrl: '...' hasValidUrl: true/false
```

### Step 3: Check imageUrl Values

Look for these potential issues:

#### Issue A: imageUrl is `undefined` or `null`
```
[ImageLayerPreview] Component mounted with imageUrl: undefined hasValidUrl: false
```
**Cause**: GraphQL query is not including `imageUrl` field
**Solution**: Check GraphQL query in backend and frontend includes `imageUrl`

#### Issue B: imageUrl is empty string `''`
```
[ImageLayerPreview] Component mounted with imageUrl: '' hasValidUrl: false
```
**Cause**: Layer was created without setting an imageUrl
**Solution**: Use ImageLayerEditor to set a valid URL, or click "Use Placeholder Image"

#### Issue C: imageUrl is set but image fails to load
```
[ImageLayerPreview] Component mounted with imageUrl: 'https://...' hasValidUrl: true
[ImageLayerPreview] Failed to load image: 'https://...'
```
**Cause**:
- Invalid URL
- CORS issues (image server doesn't allow cross-origin requests)
- Network error
- Image doesn't exist at that URL

**Solution**:
- Verify URL is correct and accessible
- Try opening URL in a new browser tab
- Check browser Network tab for CORS errors
- Use a different image host that allows cross-origin requests

### Step 4: Test Saving

1. Set an imageUrl in the ImageLayerEditor
2. Click Save
3. Watch console for:

```
[edit.vue] Saving overlay project: {...}
[edit.vue] Converting layer to API format: { type: 'IMAGE', imageUrl: '...', ... }
[edit.vue] Layers to be saved: [...]
```

4. Verify `imageUrl` is present in the data being sent to API

### Step 5: Test Loading After Save

1. Reload the page or navigate away and back
2. Watch console for the loading flow again
3. Verify `imageUrl` is loaded from API

```
[edit.vue] Loading overlay data from API: {...}
[edit.vue] Converting layer from API: { imageUrl: '...' }
```

## Common Problems & Solutions

### Problem 1: "No image URL" shown in editor

**Symptom**: The placeholder icon and "No image URL" text is displayed

**Cause**: `imageUrl` is empty, null, or undefined

**Solution**:
1. Open the layer properties panel
2. Enter a valid image URL in the "Image URL" field
3. Or click "Use Placeholder Image" button

### Problem 2: "Failed to load image" shown

**Symptom**: Error icon and "Failed to load image" text is displayed

**Causes**:
- Invalid URL format
- CORS policy blocking the image
- Image doesn't exist at the URL
- Network connectivity issues

**Solutions**:
1. Verify URL is correct (try opening in new tab)
2. Use images from CORS-friendly hosts like:
   - `https://via.placeholder.com/`
   - `https://picsum.photos/`
   - Your own server with CORS enabled
3. Check browser DevTools Network tab for specific error

### Problem 3: Image shows in editor but not in overlay

**Symptom**: ImageLayerPreview shows image correctly, but overlay page doesn't display it

**Cause**: Different rendering logic or data not synced

**Solution**:
1. Check `frontend/overlays/src/components/image-layer.vue`
2. Verify `use-overlays.ts` and `use-custom-overlay.ts` include `imageUrl` in GraphQL query
3. Check browser console on overlay page for errors

### Problem 4: imageUrl not saved to database

**Symptom**: After saving, reloading shows empty imageUrl

**Cause**: Backend is not saving imageUrl field

**Solution**:
1. Check GraphQL mutation includes `imageUrl` in `ChannelOverlayLayerSettingsInput`
2. Verify backend resolver properly saves `imageUrl` to database
3. Check database migration includes `image_url` column

## Removing Debug Logging

Once the issue is resolved, remove console.log statements from:

1. `frontend/dashboard/src/features/overlay-builder/OverlayBuilder.vue` (lines with `console.log`)
2. `frontend/dashboard/src/features/overlay-builder/components/Canvas.vue` (`@vue:mounted` handler)
3. `frontend/dashboard/src/features/overlay-builder/components/ImageLayerPreview.vue` (all console statements)
4. `frontend/dashboard/src/components/registry/overlays/edit.vue` (all console statements)

Search for `console.log` and `console.error` in these files and remove debugging statements.

## Testing Checklist

- [ ] Create new IMAGE layer - shows placeholder image
- [ ] Enter custom image URL - image loads and displays
- [ ] Enter invalid URL - shows "Failed to load image" error
- [ ] Save overlay with IMAGE layer - success message shown
- [ ] Reload page - IMAGE layer loads with correct imageUrl
- [ ] Edit imageUrl - updates in real-time
- [ ] View overlay on `/overlays` page - IMAGE layer displays correctly
- [ ] Test with CORS-enabled image - works
- [ ] Test with CORS-blocked image - shows error gracefully

## Next Steps

If issue persists after following this guide:

1. Share console logs from the debug flow
2. Check Network tab in DevTools for GraphQL requests/responses
3. Verify backend logs for any errors
4. Check database directly to see if `image_url` is being stored

## Related Files

### Frontend (Dashboard)
- `frontend/dashboard/src/features/overlay-builder/OverlayBuilder.vue`
- `frontend/dashboard/src/features/overlay-builder/components/Canvas.vue`
- `frontend/dashboard/src/features/overlay-builder/components/ImageLayerPreview.vue`
- `frontend/dashboard/src/features/overlay-builder/components/layer-editors/ImageLayerEditor.vue`
- `frontend/dashboard/src/features/overlay-builder/composables/useOverlayBuilder.ts`
- `frontend/dashboard/src/components/registry/overlays/edit.vue`

### Frontend (Overlays)
- `frontend/overlays/src/pages/overlays.vue`
- `frontend/overlays/src/components/image-layer.vue`
- `frontend/overlays/src/composables/overlays/use-overlays.ts`
- `frontend/overlays/src/composables/overlays/use-custom-overlay.ts`

### Backend (GraphQL)
- `apps/api-gql/internal/delivery/gql/schema/overlays/overlays-custom.graphql`

### Database
- Migration files for `channel_overlay_layers` table
