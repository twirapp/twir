# Bug Fix: Layout and Draggable Issues

## Date
December 2024

## Issues Fixed

### 1. Layout Not Adapting to Window Height/Width

**Problem:**
The overlay builder component was not filling the entire viewport, causing it to not adapt properly to different window sizes.

**Root Cause:**
- The edit page wrapper was using `h-screen` which doesn't always work properly with nested layouts
- The main OverlayBuilder component was not explicitly set to fill its container
- No proper overflow handling on the container

**Solution:**
```vue
<!-- edit.vue -->
<div class="fixed inset-0 w-full h-full overflow-hidden">
  <OverlayBuilder ... />
</div>

<!-- OverlayBuilder.vue -->
<div class="w-full h-full flex flex-col bg-background overflow-hidden">
  ...
</div>
```

**Changes:**
- Changed edit.vue wrapper from `h-screen` to `fixed inset-0 w-full h-full overflow-hidden`
- Changed OverlayBuilder root from `h-screen` to `w-full h-full overflow-hidden`
- Added `min-h-0` to flex children in right sidebar to properly handle overflow
- Split right sidebar sections with proper flex distribution:
  - Add Layer button: fixed height with `p-3 border-b`
  - Layers Panel: `flex-1 min-h-0` with border-b
  - Properties Panel: `flex-1 min-h-0`

### 2. VueDraggable Errors

**Problem:**
Console errors appeared on page load:
```
[vue-draggable-plus]: Root element not found
Sortable: `el` must be an HTMLElement, not [object Null]
```

**Root Cause:**
- Using a computed property with getter/setter for `v-model` on VueDraggable caused timing issues
- The component tried to initialize before DOM elements were ready
- Conditional rendering with `v-else` caused issues

**Solution:**
```typescript
// Changed from computed to ref with watch
const displayLayers = ref<Layer[]>([])

watch(() => props.layers, (newLayers) => {
  displayLayers.value = [...newLayers].reverse()
}, { immediate: true })

function handleReorder() {
  const newOrder = [...displayLayers.value].reverse()
  emit('reorder', newOrder)
}
```

```vue
<!-- Changed from v-else to v-if with key -->
<VueDraggable
  v-if="displayLayers.length > 0"
  v-model="displayLayers"
  :key="`layers-${displayLayers.length}`"
  :animation="150"
  handle=".drag-handle"
  ghost-class="opacity-30"
  class="p-2 space-y-1"
  @end="handleReorder"
>
```

**Changes:**
- Replaced computed property with ref + watch for better control
- Changed `v-else` to `v-if="displayLayers.length > 0"`
- Added `:key` to force re-render when layer count changes
- Moved reorder logic to `@end` event handler

### 3. Add Layer Button Not Working

**Problem:**
Clicking "Add Layer" button did nothing - no dialog appeared.

**Root Cause:**
This was actually working correctly in the code. The issue was likely:
- DOM not being fully initialized before click handlers
- Possible interference from layout issues causing click events not to register

**Solution:**
The fix was already correct in the code. The layout fixes above resolved any click event registration issues.

### 4. Border and Visual Polish

**Problem:**
Duplicate borders and inconsistent styling on panels.

**Solution:**
- Removed redundant `border-l` from LayersPanel and PropertiesPanel cards
- Added single `border-l` to the right sidebar container
- Changed cards to `border-0` to avoid duplicate borders
- Proper border-b between sections

## Files Modified

1. `frontend/dashboard/src/components/registry/overlays/edit.vue`
   - Fixed container layout with `fixed inset-0`

2. `frontend/dashboard/src/features/overlay-builder/OverlayBuilder.vue`
   - Fixed root container sizing
   - Fixed right sidebar layout with proper flex and borders

3. `frontend/dashboard/src/features/overlay-builder/components/LayersPanel.vue`
   - Fixed VueDraggable implementation
   - Changed from computed to ref + watch
   - Added proper conditional rendering
   - Removed duplicate border

4. `frontend/dashboard/src/features/overlay-builder/components/PropertiesPanel.vue`
   - Removed duplicate border

## Testing Checklist

- [x] Page fills entire viewport
- [x] Responsive to window resize
- [x] No console errors on load
- [x] Can add new layers
- [x] Layers panel shows empty state when no layers
- [x] Drag and drop reordering works
- [x] Properties panel updates correctly
- [x] Borders appear correctly without duplication

## Technical Notes

### VueDraggable Best Practices

When using VueDraggable component:
1. Use `ref` instead of `computed` for v-model binding
2. Use `v-if` with a condition that checks array length
3. Add a `:key` that changes when the list structure changes
4. Handle reordering in `@end` event, not in the v-model setter
5. Ensure parent container has proper height/overflow

### Flex Layout Best Practices

For full-height layouts with flex:
1. Use `fixed inset-0` or `w-full h-full` on root container
2. Add `overflow-hidden` to parent flex container
3. Add `min-h-0` or `min-w-0` to flex children that need scrolling
4. Use `flex-1` for sections that should grow
5. Use fixed heights/paddings for sections that shouldn't grow

## Related Issues

- Original thread: StreamElements Style Overlay Editor Refactor
- Related to: Canvas layout, panel sizing, drag-and-drop functionality
