# Migration Guide: Integrating the New Overlay Builder

This document provides step-by-step instructions for integrating the new modern overlay builder into the existing Twir dashboard.

## Overview

The new overlay builder is a complete rewrite with:
- Modern StreamElements-like UI/UX
- Better layer management with drag-and-drop
- Full keyboard shortcuts support
- Undo/redo functionality
- Alignment guides and snapping
- Improved state management

## What's Been Created

### New Files Structure
```
frontend/dashboard/src/features/overlay-builder/
├── types/
│   └── index.ts                    # TypeScript types and interfaces
├── composables/
│   └── useOverlayBuilder.ts        # Main state management composable
├── components/
│   ├── BuilderToolbar.vue          # Top toolbar with actions
│   ├── LayersPanel.vue             # Layers management panel
│   ├── PropertiesPanel.vue         # Properties editor panel
│   └── Canvas.vue                  # Main canvas component (incomplete)
├── OverlayBuilder.vue              # Main builder component
├── README.md                       # Documentation
└── MIGRATION.md                    # This file
```

## Installation & Setup

### 1. Install Required Dependencies

Already installed:
- ✅ `nanoid` - For unique ID generation

Existing dependencies used:
- ✅ `vue3-moveable` - For drag/resize/rotate
- ✅ `vue-draggable-plus` - For layer reordering
- ✅ `lucide-vue-next` - For icons

### 2. Complete the Canvas Component

The `Canvas.vue` component was partially created but needs to be completed. Create the file:

**Path:** `frontend/dashboard/src/features/overlay-builder/components/Canvas.vue`

Add the Moveable component at the end of the template (after line 258):

```vue
<Moveable
  v-if="selectedLayerIds.length > 0 && selectedLayers.every(l => !l.locked)"
  :target="moveableTargets"
  :draggable="true"
  :resizable="true"
  :rotatable="true"
  :snappable="snapToGrid"
  :snapThreshold="5"
  :bounds="{ left: 0, top: 0, right: canvasWidth, bottom: canvasHeight }"
  :origin="false"
  :renderDirections="['nw', 'n', 'ne', 'w', 'e', 'sw', 's', 'se']"
  @drag="onDrag"
  @drag-end="onDragEnd"
  @resize="onResize"
  @rotate="onRotate"
/>
```

Close the template tags:
```vue
      </div>
    </div>
  </div>
</template>
```

### 3. Create Missing Shadcn Components

Check if these shadcn-vue components exist. If not, install them:

```bash
cd frontend/dashboard
bun run shadcn add dialog
bun run shadcn add separator
bun run shadcn add tooltip
bun run shadcn add slider
bun run shadcn add tabs
bun run shadcn add scroll-area
```

### 4. Add Missing i18n Keys

Add these translation keys to your i18n localization files:

```typescript
// en.json or equivalent
{
  "overlaysRegistry": {
    "layers": "Layers",
    "properties": "Properties",
    // ... existing keys
  }
}
```

### 5. Create HTML Layer Renderer Component

Create a component to render HTML layers with live updates:

**Path:** `frontend/dashboard/src/features/overlay-builder/components/HtmlLayerRenderer.vue`

```vue
<script setup lang="ts">
import { useIntervalFn } from '@vueuse/core'
import { transform as transformNested } from 'nested-css-to-flat'
import { computed, nextTick, ref, watch } from 'vue'
import { useOverlaysParseHtml } from '@/api/registry'
import type { Layer } from '../types'

const props = defineProps<{
  layer: Layer
}>()

const fetcher = useOverlaysParseHtml()
const exampleValue = ref('')

const { pause, resume } = useIntervalFn(
  async () => {
    const html = props.layer.settings.htmlOverlayHtml ?? ''
    const data = await fetcher.mutateAsync(html)
    exampleValue.value = data ?? ''
  },
  (props.layer.settings.htmlOverlayDataPollSecondsInterval ?? 5) * 1000,
  { immediate: true, immediateCallback: true }
)

const executeFunc = computed(() => {
  const js = props.layer.settings.htmlOverlayJs ?? ''
  return new Function(`${js}; onDataUpdate();`)
})

watch(exampleValue, async () => {
  await nextTick()
  executeFunc.value?.()
})

watch(
  () => props.layer.periodicallyRefetchData,
  (shouldRefetch) => {
    if (!shouldRefetch) pause()
    else resume()
  },
  { immediate: true }
)

const transformedCss = computed(() => {
  const css = props.layer.settings.htmlOverlayCss ?? ''
  return transformNested(`#layer-content-${props.layer.id} { ${css} }`)
})
</script>

<template>
  <div class="w-full h-full">
    <component :is="'style'">
      {{ transformedCss }}
    </component>
    <div
      :id="`layer-content-${layer.id}`"
      class="w-full h-full"
      v-html="exampleValue"
    />
  </div>
</template>
```

### 6. Create Code Editor Component

Create a Monaco-based code editor for HTML/CSS/JS:

**Path:** `frontend/dashboard/src/features/overlay-builder/components/CodeEditor.vue`

```vue
<script setup lang="ts">
import { NAlert, NFormItem, NInputNumber, NSelect, NSwitch, NTabPane, NTabs } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { useCommandsApi } from '@/api/commands/commands'
import { useKeywordsApi } from '@/api/keywords'
import { useVariablesApi } from '@/api/variables'
import { copyToClipBoard } from '@/helpers/index.js'
import { useMessage } from 'naive-ui'
import { computed } from 'vue'
import type { Layer } from '../types'

const props = defineProps<{
  layer: Layer | null
}>()

const emit = defineEmits<{
  update: [updates: Partial<Layer>]
}>()

const { t } = useI18n()
const messages = useMessage()

const html = computed({
  get: () => props.layer?.settings?.htmlOverlayHtml ?? '',
  set: (value: string) => {
    if (!props.layer) return
    emit('update', {
      settings: {
        ...props.layer.settings,
        htmlOverlayHtml: value,
      },
    })
  },
})

const css = computed({
  get: () => props.layer?.settings?.htmlOverlayCss ?? '',
  set: (value: string) => {
    if (!props.layer) return
    emit('update', {
      settings: {
        ...props.layer.settings,
        htmlOverlayCss: value,
      },
    })
  },
})

const js = computed({
  get: () => props.layer?.settings?.htmlOverlayJs ?? '',
  set: (value: string) => {
    if (!props.layer) return
    emit('update', {
      settings: {
        ...props.layer.settings,
        htmlOverlayJs: value,
      },
    })
  },
})

const pollInterval = computed({
  get: () => props.layer?.settings?.htmlOverlayDataPollSecondsInterval ?? 5,
  set: (value: number) => {
    if (!props.layer) return
    emit('update', {
      settings: {
        ...props.layer.settings,
        htmlOverlayDataPollSecondsInterval: value,
      },
    })
  },
})

const periodicallyRefetchData = computed({
  get: () => props.layer?.periodicallyRefetchData ?? true,
  set: (value: boolean) => emit('update', { periodicallyRefetchData: value }),
})

const { allVariables } = useVariablesApi()
const keywordsManager = useKeywordsApi()
const { data: keywords } = keywordsManager.useQueryKeywords()
const commandsManager = useCommandsApi()
const { data: commands } = commandsManager.useQueryCommands()

const variables = computed(() => {
  const k = keywords.value?.keywords ?? []
  const cmds = commands.value?.commands ?? []

  return [
    ...allVariables.value
      .filter(v => v.canBeUsedInRegistry)
      .map(v => {
        const name = `$(${v.isBuiltIn ? v.name : `customvar|${v.name}`})`
        return {
          label: `${name} - ${v.description || 'Your custom variable'}`,
          value: name,
        }
      }),
    ...k.map(k => ({
      label: `$(keywords.counter|${k.id}) - How many times "${k.text}" was used`,
      value: `$(keywords.counter|${k.id})`,
    })),
    ...cmds.map(c => ({
      label: `$(command.counter.fromother|${c.name}) - How many times "${c.name}" was used`,
      value: `$(command.counter.fromother|${c.name})`,
    })),
  ]
})

async function copyVariable(v: string) {
  await copyToClipBoard(v)
  messages
