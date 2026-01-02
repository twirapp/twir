# Quick Start Guide - New Overlay Builder

## 🚀 Getting Started

The new overlay builder has been integrated and is ready to use! Here's how to test it:

## ✅ Prerequisites

All required dependencies and components are already installed:
- ✅ `nanoid` - Installed
- ✅ `vue3-moveable` - Already in project
- ✅ `vue-draggable-plus` - Already in project
- ✅ All Shadcn components - Already installed

## 🧪 Testing the Builder

### 1. Navigate to Overlays
1. Open your browser and go to `http://localhost:5173/dashboard/overlays` (or your dev URL)
2. Click on an existing overlay OR click "Create New Overlay"
3. You should now see the **new modern builder** instead of the old editor!

### 2. What You Should See

The new builder has a professional layout with:

**Top Bar:**
- Save button
- Undo/Redo buttons
- Copy/Cut/Paste/Delete buttons
- Alignment tools
- Zoom controls
- Grid toggle

**Center Canvas:**
- Dark background with grid (toggleable)
- Your overlay layers displayed
- Blue alignment guides when dragging

**Right Sidebar - Top Half:**
- "Add Layer" button
- Layers panel showing all layers
- Drag layers to reorder
- Eye icon to toggle visibility
- Lock icon to prevent editing
- Duplicate and delete buttons

**Right Sidebar - Bottom Half:**
- Properties panel for selected layer
- Position (X, Y)
- Size (Width, Height)
- Rotation slider
- Opacity slider
- Visibility and Lock toggles

### 3. Basic Operations to Test

#### Adding Layers
1. Click "Add Layer" button (top of right sidebar)
2. Select "HTML Layer"
3. Layer appears on canvas

#### Moving Layers
1. Click on a layer in the canvas
2. Drag it around
3. Notice blue alignment guides appear when near other layers

#### Resizing Layers
1. Click on a layer
2. Drag the corner handles to resize
3. Hold Shift for proportional resize

#### Rotating Layers
1. Select a layer
2. Use the rotation slider in Properties panel
3. Or drag the rotation handle (circular icon)

#### Multi-Selection
1. Click first layer
2. Hold Ctrl (or Cmd on Mac)
3. Click another layer
4. Both are now selected!

#### Alignment
1. Select 2 or more layers (Ctrl+Click)
2. Click alignment buttons in toolbar
3. Layers align to left/center/right/top/middle/bottom

#### Undo/Redo
1. Make any change
2. Press Ctrl+Z to undo
3. Press Ctrl+Y to redo

### 4. Keyboard Shortcuts to Test

| Shortcut | Action |
|----------|--------|
| `Ctrl+S` | Save overlay |
| `Ctrl+Z` | Undo |
| `Ctrl+Y` | Redo |
| `Ctrl+C` | Copy selected layers |
| `Ctrl+V` | Paste layers |
| `Ctrl+X` | Cut layers |
| `Ctrl+D` | Duplicate layers |
| `Delete` | Delete selected layers |
| `Ctrl+A` | Select all layers |

### 5. Layer Management to Test

In the Layers Panel (right sidebar):

1. **Reorder:** Drag a layer up/down in the list
2. **Visibility:** Click eye icon to hide/show layer
3. **Lock:** Click lock icon to prevent editing
4. **Duplicate:** Click copy icon to duplicate
5. **Delete:** Click trash icon to delete
6. **Select:** Click layer name to select it

### 6. Properties Panel to Test

When you select a single layer:

1. **Name:** Edit layer name
2. **Position Tab:**
   - Change X, Y position manually
   - Change Width, Height
   - Adjust Rotation slider (0-360°)
3. **Appearance Tab:**
   - Adjust Opacity slider (0-100%)
   - Toggle Visibility switch
   - Toggle Locked switch
4. **HTML Settings (for HTML layers):**
   - Click "Edit Code" button
   - Toggle Auto Refresh
   - Set Update Interval

## 🐛 Known Limitations

### What's Working
- ✅ All core functionality
- ✅ Layer management (add, remove, reorder)
- ✅ Transform controls (move, resize, rotate)
- ✅ Multi-selection
- ✅ Alignment tools
- ✅ Undo/Redo
- ✅ Keyboard shortcuts
- ✅ Save/Load overlays

### What Needs Enhancement
- ⚠️ HTML layer rendering - Shows placeholder text instead of live HTML preview
- ⚠️ Code Editor - Opens placeholder modal (Monaco editor integration pending)

### To Fix These:
See `MIGRATION.md` steps 4-7 for creating HtmlLayerRenderer and CodeEditor components.

## 🎯 Expected Behavior

### When Creating New Overlay
1. Click "Create New" in overlays list
2. Builder opens with empty canvas (1920x1080)
3. Click "Add Layer" to add your first layer
4. Configure and position layers
5. Click Save (or Ctrl+S)
6. Overlay is created and you can use it!

### When Editing Existing Overlay
1. Click an existing overlay in the list
2. Builder opens with all existing layers loaded
3. Layers appear at their saved positions
4. Make changes
5. Click Save to update

### When Using Keyboard Shortcuts
1. All shortcuts work when builder has focus
2. Ctrl+Z/Y undo/redo any action
3. Ctrl+C/V copies/pastes layers (with +20px offset)
4. Delete key removes selected layers
5. Ctrl+S saves without confirmation

## 🎨 UI Features

### Visual Feedback
- **Selected layers:** Blue border and label
- **Hover:** Gray border appears
- **Locked layers:** Gray border, lock icon, can't edit
- **Hidden layers:** Grayed out in layers panel
- **Alignment guides:** Blue lines when dragging near edges

### Grid System
- Click grid icon in toolbar to toggle grid display
- Click snap icon to enable/disable snap-to-grid
- Grid helps align layers visually
- Snap makes layers jump to grid positions

### Zoom Controls
- Click `-` to zoom out
- Click `+` to zoom in
- Click percentage to reset to 100%
- Zoom affects view, not actual layer sizes

## 🔍 Debugging

### If Builder Doesn't Load
1. Check browser console for errors
2. Verify you're on `/dashboard/registry/overlays/:id` route
3. Check that overlay data is loading
4. Look for missing Shadcn components

### If Layers Don't Appear
1. Check if overlay has layers in database
2. Verify layer data format matches expectations
3. Check canvas dimensions (should be 1920x1080 default)

### If Controls Don't Work
1. Ensure layer is selected (blue border)
2. Check if layer is locked (unlock it)
3. Try clicking canvas first to give it focus
4. Check browser console for JavaScript errors

### Common Issues

**"Nothing happens when I click Add Layer"**
- Check browser console for errors
- Verify Dialog component is imported correctly

**"Layers overlap incorrectly"**
- Check zIndex values in layers panel
- Reorder layers by dragging in the panel

**"Can't select layer"**
- Layer might be locked (click lock icon to unlock)
- Layer might be hidden (click eye icon to show)
- Try clicking directly on layer content

**"Undo doesn't work"**
- Make sure you made a change first
- Check that history isn't empty (max 50 states)
- Try Ctrl+Z (not just Z)

## 📝 Next Steps

### For Development
1. ✅ Test basic functionality - Make sure everything works
2. ⚠️ Add HTML live preview - See MIGRATION.md step 4
3. ⚠️ Add Monaco code editor - See MIGRATION.md step 5
4. 📋 Add i18n translations - Add missing keys
5. 🎨 Customize styling - Adjust colors/spacing if needed

### For Production
1. Test with real overlays
2. Test with maximum 15 layers
3. Test undo/redo extensively
4. Test on different screen sizes
5. Test keyboard shortcuts thoroughly
6. Verify backwards compatibility
7. Test save/load multiple times

## 🎉 Success!

If you can:
- ✅ Create a new overlay
- ✅ Add layers
- ✅ Move and resize layers
- ✅ Use keyboard shortcuts
- ✅ Save and reload

**Congratulations!** The new overlay builder is working! 🎊

## 📚 More Information

- Full documentation: `README.md`
- Integration guide: `MIGRATION.md`
- Complete checklist: `CHECKLIST.md`
- Feature summary: `SUMMARY.md`

## 🆘 Need Help?

1. Check browser console for errors
2. Review documentation files
3. Look at example-integration.vue
4. Verify all imports are correct
5. Check that route meta has `fullScreen: true`
