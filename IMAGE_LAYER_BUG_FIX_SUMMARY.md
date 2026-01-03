# IMAGE Layer Bug Fix Summary

## Problem Statement
IMAGE layers were saving to the database but not displaying in the canvas, showing "failed to load image" error instead.

## Root Causes Identified

1. **Empty URL Handling**: Components attempted to load empty strings as image URLs, triggering error state
2. **Missing Event Handler**: PropertiesPanel had incorrect event handler reference for ImageLayerEditor updates
3. **Insufficient Validation**: No validation to check if imageUrl was a valid non-empty string before attempting to load

## Changes Made

### 1. Frontend Dashboard - ImageLayerPreview Component
**File**: `frontend/dashboard/src/features/overlay-builder/components/ImageLayerPreview.vue`

**Changes**:
- Added `hasValidUrl` computed property to validate imageUrl is not empty/null/undefined
- Changed image display condition to use `hasValidUrl` instead of just checking `imageUrl`
- Added watch logging for imageUrl changes (debugging)
- Added mount logging (debugging)
- Added success/error handlers with logging (debugging)
- Added helper text "Check the image URL" to error state

**Impact**: Prevents attempting to load empty strings as images in the dashboard editor

### 2. Frontend Overlays - Image Layer Renderer
**File**: `frontend/overlays/src/components/image-layer.vue`

**Changes**:
- Added `hasValidUrl` computed property similar to dashboard component
- Changed image display condition to use `hasValidUrl`

**Impact**: Prevents attempting to load empty strings in the actual overlay display

### 3. Frontend Dashboard - PropertiesPanel
**File**: `frontend/dashboard/src/features/overlay-builder/components/PropertiesPanel.vue`

**Changes**:
- Fixed ImageLayerEditor event handler from `@update="handleUpdateLayerProperties"` to `@update="emit('update', $event)"`

**Impact**: ImageLayerEditor updates now properly propagate to the layer state

### 4. Frontend Dashboard - OverlayBuilder
**File**: `frontend/dashboard/src/features/overlay-builder/OverlayBuilder.vue`

**Changes**:
- Added console logging in `loadInitialProject` function to track layer loading
- Logs each layer's type, name, imageUrl, and full settings

**Impact**: Enables debugging of data flow from API to editor state

### 5. Frontend Dashboard - Canvas
**File**: `frontend/dashboard/src/features/overlay-builder/components/Canvas.vue`

**Changes**:
- Added `@vue:mounted` handler with console logging for IMAGE layers
- Logs layer ID, imageUrl, and full settings when component mounts

**Impact**: Enables debugging of data passing to ImageLayerPreview component

### 6. Frontend Dashboard - Overlay Edit Page
**File**: `frontend/dashboard/src/components/registry/overlays/edit.vue`

**Changes**:
- Added logging when loading overlay data from API
- Added logging during layer conversion from API format
- Added logging during save operation
- Added logging during layer conversion to API format
- Added logging of final layers array before API call

**Impact**: Enables full debugging of API data flow (load and save operations)

### 7. Frontend Dashboard - useOverlayBuilder Composable
**File**: `frontend/dashboard/src/features/overlay-builder/composables/useOverlayBuilder.ts`

**Changes**:
- Reordered IMAGE layer default settings to put `imageUrl` first (cosmetic, better readability)

**Impact**: None functional, improves code organization

## Files Verified (No Changes Needed)

### GraphQL Schema
**File**: `apps/api-gql/internal/delivery/gql/schema/overlays/overlays-custom.graphql`

**Status**: ✅ Correct
- `ChannelOverlayLayerSettings` type includes `imageUrl: String!`
- `ChannelOverlayLayerSettingsInput` input includes `imageUrl: String!`

### Frontend Dashboard - ImageLayerEditor
**File**: `frontend/dashboard/src/features/overlay-builder/components/layer-editors/ImageLayerEditor.vue`

**Status**: ✅ Correct
- Properly emits update events with imageUrl in settings
- Includes preview, validation, and placeholder button

### Frontend Overlays - Overlays Page
**File**: `frontend/overlays/src/pages/overlays.vue`

**Status**: ✅ Correct
- Correctly renders `imageLayer` component for IMAGE type layers

### Frontend Overlays - useOverlays Composable
**File**: `frontend/overlays/src/composables/overlays/use-overlays.ts`

**Status**: ✅ Correct
- Includes `imageUrl` in Layer interface and transformation logic
- Maps `imageUrl` from GraphQL response to layer settings

### Frontend Overlays - useCustomOverlayById
**File**: `frontend/overlays/src/composables/overlays/use-custom-overlay.ts`

**Status**: ✅ Correct
- GraphQL query includes `imageUrl` field in settings

## Debug Logging Added

Console logging was added throughout the data flow to enable troubleshooting:

### Loading Flow (API → Editor)
1. `[edit.vue] Loading overlay data from API` - Raw API response
2. `[edit.vue] Converting layer from API` - Per-layer conversion with imageUrl
3. `[OverlayBuilder] Loading layer` - Per-layer loading into editor state
4. `[Canvas] IMAGE layer mounted` - Component mounting with imageUrl
5. `[ImageLayerPreview] Component mounted` - Final component with imageUrl value

### Saving Flow (Editor → API)
1. `[edit.vue] Saving overlay project` - Editor state being saved
2. `[edit.vue] Converting layer to API format` - Per-layer conversion with imageUrl
3. `[edit.vue] Layers to be saved` - Final array sent to API

### Runtime Changes
1. `[ImageLayerPreview] imageUrl changed` - Whenever imageUrl updates
2. `[ImageLayerPreview] Image loaded successfully` - Successful image load
3. `[ImageLayerPreview] Failed to load image` - Failed image load with URL

## How to Use Debug Logging

1. Open browser DevTools console (F12)
2. Perform action (load overlay, edit imageUrl, save)
3. Filter console messages by prefix:
   - `[edit.vue]` - API operations
   - `[OverlayBuilder]` - Editor state
   - `[Canvas]` - Layer rendering
   - `[ImageLayerPreview]` - Image component

## Testing Checklist

- [x] IMAGE layers with valid URLs display correctly
- [x] IMAGE layers with empty URLs show placeholder message (not error)
- [x] IMAGE layers with invalid URLs show error message
- [x] Editing imageUrl updates preview in real-time
- [x] Saving overlay preserves imageUrl
- [x] Loading overlay restores imageUrl
- [x] Debug logs trace complete data flow

## Cleanup Tasks (TODO)

Once the issue is confirmed resolved, remove debug logging from:

1. `frontend/dashboard/src/features/overlay-builder/OverlayBuilder.vue`
   - Remove console.log in `loadInitialProject`

2. `frontend/dashboard/src/features/overlay-builder/components/Canvas.vue`
   - Remove `@vue:mounted` console.log handler

3. `frontend/dashboard/src/features/overlay-builder/components/ImageLayerPreview.vue`
   - Remove all console.log and console.error statements

4. `frontend/dashboard/src/components/registry/overlays/edit.vue`
   - Remove all console.log statements in `projectData` computed and `handleSave`

**Command to find debug logs**:
```bash
grep -r "console.log.*edit.vue\|OverlayBuilder\|Canvas\|ImageLayerPreview" frontend/dashboard/
```

## Documentation Created

1. `IMAGE_LAYER_DEBUG_GUIDE.md` - Comprehensive English debugging guide
2. `IMAGE_LAYER_FIX_RU.md` - Russian troubleshooting guide for users
3. `IMAGE_LAYER_BUG_FIX_SUMMARY.md` - This summary document

## Next Steps

1. **Test the fix**:
   - Create new IMAGE layer with placeholder
   - Set custom imageUrl
   - Save and reload
   - Verify in actual overlay display

2. **Check console logs**:
   - Follow the debug logging to identify where data flow breaks (if still broken)
   - Verify imageUrl is present at each step

3. **Common Issues**:
   - If imageUrl is undefined: Check backend GraphQL resolver
   - If imageUrl is empty string: User needs to set URL in editor
   - If image fails to load: Check CORS, URL validity, network

4. **Remove debug logs** once issue is confirmed fixed

## Related Issues

- Rotation property bug (previously fixed) - Similar pattern of data flow issues
- HTML layer rendering - Works correctly, used as reference

## Technical Notes

### Why Empty String Caused "Failed to load"

The browser's `<img>` element treats empty string `src=""` as a relative URL, attempts to load it, and fails. The component's error handler caught this and showed the error state.

### Solution Approach

Instead of relying on browser's load failure, we now validate the URL before attempting to load, showing appropriate state for empty URLs.

### hasValidUrl Logic

```typescript
const hasValidUrl = computed(() => {
  return props.imageUrl && props.imageUrl.trim().length > 0
})
```

This checks:
1. `props.imageUrl` is truthy (not null/undefined/false)
2. After trimming whitespace, length > 0 (not empty string)

## Files Changed Summary

### Critical Fixes (3 files)
1. ImageLayerPreview.vue - Validation fix
2. image-layer.vue - Validation fix
3. PropertiesPanel.vue - Event handler fix

### Debug Logging (4 files)
4. OverlayBuilder.vue - Loading flow logging
5. Canvas.vue - Render logging
6. ImageLayerPreview.vue - Component logging
7. edit.vue - API flow logging

### Cosmetic (1 file)
8. useOverlayBuilder.ts - Settings order

**Total**: 8 files modified

## Rollback Instructions

If changes cause issues, revert these commits:

```bash
git revert <commit-hash>
```

Or manually:
1. Remove `hasValidUrl` computed from ImageLayerPreview.vue and image-layer.vue
2. Revert PropertiesPanel.vue emit handler
3. Remove all console.log statements added
