# Overlay Builder Implementation Checklist

## ✅ Completed

- [x] Created type definitions (`types/index.ts`)
- [x] Created main state management composable (`composables/useOverlayBuilder.ts`)
- [x] Created BuilderToolbar component with all actions
- [x] Created LayersPanel component with drag-to-reorder
- [x] Created PropertiesPanel component with all controls
- [x] Created Canvas component with Moveable integration
- [x] Created main OverlayBuilder component
- [x] Fixed TypeScript errors in Canvas and OverlayBuilder
- [x] Installed `nanoid` dependency
- [x] Created comprehensive documentation (README.md, MIGRATION.md, SUMMARY.md)
- [x] Created example integration file

## 🔧 Required to Complete

### 1. Install Shadcn Components
Run these commands in `frontend/dashboard`:

```bash
cd frontend/dashboard

# Install required shadcn components
bun run shadcn add dialog
bun run shadcn add separator
bun run shadcn add tooltip
bun run shadcn add slider
bun run shadcn add tabs
bun run shadcn add scroll-area
```

**Priority:** HIGH - Builder won't work without these components

### 2. Update Existing Overlay Edit Page
Replace content in `components/registry/overlays/edit.vue` with the new builder integration.

**File:** `twir/frontend/dashboard/src/components/registry/overlays/edit.vue`

See `example-integration.vue` for reference implementation.

**Priority:** HIGH - This is the main integration point

### 3. Add i18n Translation Keys
Add missing translation keys for layers and properties panels.

**Files to update:**
- `frontend/dashboard/src/locales/en.json` (or equivalent)
- `frontend/dashboard/src/locales/ru.json` (or equivalent)

```json
{
  "overlaysRegistry": {
    "layers": "Layers",
    "properties": "Properties",
    "validations": {
      "name": "Name is required and must be less than 30 characters",
      "layers": "Overlay must have between 1 and 15 layers"
    }
  }
}
```

**Priority:** MEDIUM - Will show missing translation keys but won't break functionality

### 4. Create HTML Layer Renderer Component
Create a component to render HTML layers with live preview in the canvas.

**File:** `twir/frontend/dashboard/src/features/overlay-builder/components/HtmlLayerRenderer.vue`

**Reference:** Use existing `components/registry/overlays/layers/html.vue` as a template

**Key features:**
- Render HTML with variable substitution
- Execute JavaScript in sandbox
- Apply CSS with scoped styles
- Poll for data updates if `periodicallyRefetchData` is enabled

**Priority:** MEDIUM - Canvas will show placeholder without this

### 5. Create Code Editor Component
Create a full-featured code editor modal for HTML/CSS/JS editing.

**File:** `twir/frontend/dashboard/src/features/overlay-builder/components/CodeEditor.vue`

**Reference:** Use existing `components/registry/overlays/layers/htmlForm.vue` as a template

**Key features:**
- Monaco editor tabs for HTML, CSS, JS
- Variables selector with copy-to-clipboard
- Update interval settings
- Auto-refresh toggle

**Priority:** MEDIUM - Can use properties panel initially, but code editor is better UX

### 6. Wire HTML Renderer into Canvas
Update `OverlayBuilder.vue` to use the HtmlLayerRenderer in the canvas slot.

```vue
<Canvas ...>
  <template #layer-content="{ layer }">
    <HtmlLayerRenderer
      v-if="layer.type === ChannelOverlayLayerType.Html"
      :layer="layer"
    />
    <div v-else class="text-white/70 text-sm">
      {{ layer.name }}
    </div>
  </template>
</Canvas>
```

**Priority:** MEDIUM - Needed for live preview

### 7. Wire Code Editor into Properties Panel
Update `OverlayBuilder.vue` to open the CodeEditor modal instead of placeholder.

```vue
<Dialog v-model:open="showCodeEditor">
  <DialogContent class="max-w-4xl h-[80vh]">
    <CodeEditor
      :layer="builder.activeLayer.value"
      @update="handleUpdateLayerProperties"
    />
  </DialogContent>
</Dialog>
```

**Priority:** LOW - Properties panel works without it

## 🧪 Testing Steps

### Basic Functionality
- [ ] Builder loads without errors
- [ ] Can create new overlay
- [ ] Can add HTML layer
- [ ] Can move layer with mouse drag
- [ ] Can resize layer with corner handles
- [ ] Can rotate layer (if enabled)

### Layer Management
- [ ] Can select single layer by clicking
- [ ] Can multi-select with Ctrl/Cmd + Click
- [ ] Can reorder layers in panel by dragging
- [ ] Can toggle layer visibility (eye icon)
- [ ] Can lock/unlock layer (lock icon)
- [ ] Can duplicate layer
- [ ] Can delete layer

### Clipboard Operations
- [ ] Copy (Ctrl+C) copies selected layers
- [ ] Cut (Ctrl+X) cuts selected layers
- [ ] Paste (Ctrl+V) pastes layers
- [ ] Duplicate (Ctrl+D) duplicates selected layers

### Alignment & Distribution
- [ ] Can align multiple layers left/center/right
- [ ] Can align multiple layers top/middle/bottom
- [ ] Can distribute 3+ layers horizontally
- [ ] Can distribute 3+ layers vertically
- [ ] Alignment guides appear when dragging near other layers

### History Management
- [ ] Undo (Ctrl+Z) reverts last action
- [ ] Redo (Ctrl+Y) redoes undone action
- [ ] History preserved across multiple operations
- [ ] Can undo/redo multiple times

### Canvas Controls
- [ ] Zoom in button increases zoom
- [ ] Zoom out button decreases zoom
- [ ] Click zoom percentage resets to 100%
- [ ] Grid toggle shows/hides grid
- [ ] Snap toggle enables/disables snap-to-grid

### Properties Panel
- [ ] Shows properties of selected layer
- [ ] Can edit layer name
- [ ] Can edit X, Y position
- [ ] Can edit width, height
- [ ] Can adjust rotation slider
- [ ] Can adjust opacity slider
- [ ] Can toggle visibility switch
- [ ] Can toggle locked switch
- [ ] HTML settings show for HTML layers

### Saving & Loading
- [ ] Can save new overlay
- [ ] Can save existing overlay
- [ ] Can load existing overlay
- [ ] All layer properties preserved on save/load
- [ ] New properties (rotation, opacity) handled gracefully

### Edge Cases
- [ ] Can't edit locked layers
- [ ] Can't move layers outside canvas bounds
- [ ] Multi-selection excludes locked layers from transform
- [ ] Hidden layers don't interfere with selection
- [ ] Deleting selected layers clears selection
- [ ] Works with 0 layers (shows empty state)
- [ ] Works with maximum 15 layers

## 📝 Optional Enhancements

### Nice to Have
- [ ] Add layer name editing inline in layers panel
- [ ] Add layer thumbnails in layers panel
- [ ] Add context menu (right-click) on layers
- [ ] Add canvas mini-map for navigation
- [ ] Add ruler guides along edges
- [ ] Add measurement tooltips while resizing
- [ ] Add color picker for backgrounds
- [ ] Add keyboard shortcut help dialog (?)

### Future Features
- [ ] Layer groups/folders
- [ ] Layer blend modes
- [ ] Layer effects (shadow, blur, etc.)
- [ ] Animation timeline
- [ ] Templates library
- [ ] Asset manager
- [ ] Export to image/video

## 🐛 Known Issues

### To Fix
- Canvas.vue has unused `canvasContainer` ref (warning only, not critical)

### To Investigate
- Test performance with 15 layers
- Test memory usage with large undo history
- Verify Monaco editor doesn't conflict with existing instances
- Check z-index stacking with Naive UI modals

## 📚 Documentation Updates Needed

- [ ] Update main README with new builder features
- [ ] Add screenshots/GIFs of new builder UI
- [ ] Document keyboard shortcuts in user docs
- [ ] Add troubleshooting section for common issues

## 🚀 Deployment Checklist

Before deploying to production:
- [ ] All HIGH priority items completed
- [ ] All tests passing
- [ ] No console errors in browser
- [ ] Tested in Chrome, Firefox, Safari
- [ ] Tested on different screen sizes
- [ ] i18n keys added for all supported languages
- [ ] Performance tested with max layers
- [ ] Existing overlays load correctly
- [ ] Can save and load new format
- [ ] Backwards compatibility verified

## 📞 Support

If you encounter issues:
1. Check browser console for errors
2. Verify all Shadcn components are installed
3. Check that Monaco editor is properly configured
4. Review type errors in IDE
5. Refer to README.md and MIGRATION.md documentation

## 🎉 Success Criteria

The integration is complete when:
- ✅ All HIGH priority items are done
- ✅ Builder loads and functions without errors
- ✅ Can create, edit, save, and load overlays
- ✅ All keyboard shortcuts work
- ✅ All toolbar buttons function correctly
- ✅ Layers can be managed (add, remove, reorder, etc.)
- ✅ Properties can be edited
- ✅ Undo/redo works reliably
- ✅ No regressions in existing overlay functionality
