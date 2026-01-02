# Overlay Builder System

A modern, StreamElements-like overlay builder for creating custom overlays with drag-and-drop functionality, layers management, and advanced editing features.

## Features

### Core Functionality
- **Layer Management**: Add, remove, duplicate, reorder layers
- **Drag & Drop**: Move and resize layers with visual feedback
- **Alignment Guides**: Smart guides that appear when aligning with other layers
- **Grid & Snapping**: Optional grid display and snap-to-grid functionality
- **Multi-Selection**: Select multiple layers with Ctrl/Cmd + Click
- **Undo/Redo**: Full history management with Ctrl+Z/Ctrl+Y
- **Clipboard**: Copy, cut, and paste layers
- **Layer Properties**: Visibility, lock, opacity, rotation
- **Keyboard Shortcuts**: Full keyboard navigation and actions

### UI Components

#### BuilderToolbar
Top toolbar with all common actions:
- Save, Undo, Redo
- Copy, Cut, Paste, Delete
- Alignment tools (left, center, right, top, middle, bottom)
- Distribution tools (horizontal, vertical)
- Zoom controls
- Grid and snap toggles

#### LayersPanel
Photoshop-like layers panel:
- Drag to reorder layers
- Toggle visibility (eye icon)
- Toggle lock (lock icon)
- Duplicate and delete actions
- Shows layer thumbnails and dimensions
- Visual indication of selected layers

#### PropertiesPanel
Context-sensitive properties editor:
- Position (X, Y)
- Size (Width, Height)
- Rotation (0-360°)
- Opacity (0-100%)
- Visibility and Lock toggles
- Layer-specific settings (HTML code editor for HTML layers)

#### Canvas
Main editing canvas:
- Zoom and pan support
- Visual grid overlay
- Alignment guides
- Selection highlights
- Moveable integration for transforms
- Layer content rendering

## Architecture

### State Management
All state is managed through the `useOverlayBuilder` composable:

```typescript
const builder = useOverlayBuilder()

// Access state
builder.project // Current project data
builder.canvasState // Canvas zoom, pan, selection
builder.selectedLayers // Currently selected layers
builder.activeLayer // Single active layer (for properties panel)

// History
builder.undo()
builder.redo()

// Layer operations
builder.addLayer(type, options)
builder.removeLayer(layerId)
builder.updateLayer(layerId, updates)
builder.duplicateLayer(layerId)

// Selection
builder.selectLayers([layerId], addToSelection)
builder.deselectAll()

// Alignment
builder.alignLayers('left' | 'center' | 'right' | 'top' | 'middle' | 'bottom')
builder.distributeLayersHorizontally()
builder.distributeLayersVertically()

// Canvas
builder.setZoom(zoom)
builder.zoomIn()
builder.zoomOut()
```

### Type System
All types are defined in `types/index.ts`:

```typescript
interface Layer {
  id: string
  type: ChannelOverlayLayerType
  name: string
  posX: number
  posY: number
  width: number
  height: number
  rotation: number
  opacity: number
  visible: boolean
  locked: boolean
  zIndex: number
  periodicallyRefetchData: boolean
  settings: LayerSettings
}

interface OverlayProject {
  id: string
  name: string
  width: number
  height: number
  layers: Layer[]
}

interface CanvasState {
  zoom: number
  panX: number
  panY: number
  selectedLayerIds: string[]
  clipboardLayers: Layer[]
  showGrid: boolean
  snapToGrid: boolean
  gridSize: number
  showRulers: boolean
  showGuides: boolean
}
```

## Integration with Existing Code

### Replacing Old Overlay Editor

The old overlay editor was in `components/registry/overlays/edit.vue`. To integrate the new builder:

1. **Import the OverlayBuilder component**:
```vue
<script setup lang="ts">
import OverlayBuilder from '@/features/overlay-builder/OverlayBuilder.vue'
</script>
```

2. **Convert existing data to the new format**:
```typescript
const overlayData = computed(() => {
  if (!overlay.value) return null

  return {
    id: overlay.value.id,
    name: overlay.value.name,
    width: overlay.value.width,
    height: overlay.value.height,
    layers: overlay.value.layers.map((layer, index) => ({
      id: `layer-${index}`, // or use existing ID if available
      type: layer.type,
      name: `Layer ${index + 1}`,
      posX: layer.posX,
      posY: layer.posY,
      width: layer.width,
      height: layer.height,
      rotation: 0,
      opacity: 1,
      visible: true,
      locked: false,
      zIndex: index,
      periodicallyRefetchData: layer.periodicallyRefetchData,
      settings: {
        htmlOverlayHtml: layer.settings.htmlOverlayHtml,
        htmlOverlayCss: layer.settings.htmlOverlayCss,
        htmlOverlayJs: layer.settings.htmlOverlayJs,
        htmlOverlayDataPollSecondsInterval: layer.settings.htmlOverlayDataPollSecondsInterval,
      },
    })),
  }
})
```

3. **Handle save callback**:
```vue
<template>
  <OverlayBuilder
    :initial-project="overlayData"
    @save="handleSave"
  />
</template>

<script setup lang="ts">
async function handleSave(project: OverlayProject) {
  const layersInput = project.layers.map(layer => ({
    type: layer.type,
    posX: layer.posX,
    posY: layer.posY,
    width: layer.width,
    height: layer.height,
    periodicallyRefetchData: layer.periodicallyRefetchData,
    settings: {
      htmlOverlayHtml: layer.settings.htmlOverlayHtml || '',
      htmlOverlayCss: layer.settings.htmlOverlayCss || '',
      htmlOverlayJs: layer.settings.htmlOverlayJs || '',
      htmlOverlayDataPollSecondsInterval: layer.settings.htmlOverlayDataPollSecondsInterval || 5,
    },
  }))

  if (project.id) {
    await updateOverlayMutation.executeMutation({
      id: project.id,
      input: {
        name: project.name,
        width: project.width,
        height: project.height,
        layers: layersInput,
      },
    })
  } else {
    await createOverlayMutation.executeMutation({
      input: {
        name: project.name,
        width: project.width,
        height: project.height,
        layers: layersInput,
      },
    })
  }
}
</script>
```

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `Ctrl/Cmd + S` | Save |
| `Ctrl/Cmd + Z` | Undo |
| `Ctrl/Cmd + Y` or `Ctrl/Cmd + Shift + Z` | Redo |
| `Ctrl/Cmd + C` | Copy selected layers |
| `Ctrl/Cmd + X` | Cut selected layers |
| `Ctrl/Cmd + V` | Paste layers |
| `Ctrl/Cmd + D` | Duplicate selected layers |
| `Delete` or `Backspace` | Delete selected layers |
| `Ctrl/Cmd + A` | Select all layers |
| `Ctrl/Cmd + Click` | Add to selection |

## Extending the Builder

### Adding New Layer Types

1. Add the layer type to GraphQL schema if needed
2. Update the `Layer` type in `types/index.ts` if custom settings are needed
3. Add a new option in the "Add Layer" dialog in `OverlayBuilder.vue`
4. Implement the layer renderer in the Canvas slot

Example:
```vue
<Canvas>
  <template #layer-content="{ layer }">
    <HtmlLayer v-if="layer.type === 'HTML'" :layer="layer" />
    <ImageLayer v-else-if="layer.type === 'IMAGE'" :layer="layer" />
    <TextLayer v-else-if="layer.type === 'TEXT'" :layer="layer" />
  </template>
</Canvas>
```

### Custom Properties Panel Sections

Add custom sections to the PropertiesPanel for specific layer types:

```vue
<!-- In PropertiesPanel.vue -->
<div v-if="layer.type === 'HTML'" class="space-y-4">
  <h4 class="text-sm font-medium">HTML Settings</h4>
  <!-- Custom HTML settings here -->
</div>
```

## Dependencies

- `vue3-moveable`: For drag, resize, and rotate functionality
- `vue-draggable-plus`: For reorderable layers panel
- `nanoid`: For generating unique layer IDs
- `lucide-vue-next`: For icons
- Shadcn Vue components: Button, Card, Dialog, Input, Label, etc.

## Performance Considerations

- **History Management**: Limited to 50 undo states to prevent memory issues
- **Layer Rendering**: Use Vue's virtual scrolling for large numbers of layers
- **Canvas Optimization**: Only render visible layers based on viewport
- **Debouncing**: Consider debouncing property updates for better performance

## Future Enhancements

- [ ] Layer groups/folders
- [ ] Layer effects (blur, shadow, etc.)
- [ ] Animations timeline
- [ ] Templates library
- [ ] Asset manager for images/fonts
- [ ] Responsive breakpoints
- [ ] Canvas rulers and measurement tools
- [ ] Layer search/filter
- [ ] Export to different formats
- [ ] Collaboration features (if needed)
