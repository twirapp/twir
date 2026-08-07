<script setup lang="ts">
import { toRef } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardTitle } from '@/components/ui/card'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { ScrollArea } from '@/components/ui/scroll-area'
import type { ChannelOverlayLayerType } from '~/gql/graphql.js'

import type { Layer } from '../types'
import { getLayerTypeMeta } from '../layer-type-meta'
import { useLayersPanel } from '../composables/useLayersPanel'

interface Props {
	layers: Layer[]
	selectedLayerIds: string[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
	select: [layerId: string, addToSelection: boolean]
	toggleVisibility: [layerId: string]
	toggleLock: [layerId: string]
	reorder: [layers: Layer[]]
	addLayer: [type: ChannelOverlayLayerType]
	updateLayerProperties: [layerId: string, updates: Partial<Layer>]
}>()

const {
	displayLayers,
	isAddPopoverOpen,
	layerTypeOptions,
	handleAddLayerType,
	handleReorder,
	handleLayerClick,
	isLayerSelected,
	renameDraft,
	renamingLayerId,
	startRename,
	commitRename,
	cancelRename,
} = useLayersPanel(
	toRef(props, 'layers'),
	toRef(props, 'selectedLayerIds'),
	(layerId, addToSelection) => emit('select', layerId, addToSelection),
	(type) => emit('addLayer', type),
	(layers) => emit('reorder', layers),
	(layerId, updates) => emit('updateLayerProperties', layerId, updates),
)
</script>

<template>
	<Card class="flex h-full flex-col border-0 p-0">
		<div class="flex items-center justify-between border-b p-2">
			<div class="flex items-center gap-2">
				<CardTitle class="text-sm font-medium">Слои</CardTitle>
				<span class="rounded-full bg-muted px-1.5 py-0.5 text-[10px] font-medium text-muted-foreground">
					{{ layers.length }}
				</span>
			</div>
			<Popover v-model:open="isAddPopoverOpen">
				<PopoverTrigger as-child>
					<Button size="sm" class="h-7 gap-1 bg-emerald-600 px-2.5 text-xs font-medium text-white hover:bg-emerald-500">
						<Icon name="lucide:plus" class="h-3.5 w-3.5" />
						Добавить
					</Button>
				</PopoverTrigger>
				<PopoverContent align="end" class="w-80 p-2">
					<div class="grid grid-cols-2 gap-1">
						<button
							v-for="option in layerTypeOptions"
							:key="option.type"
							type="button"
							class="hover:bg-accent flex flex-col items-start gap-1 rounded-md p-2 text-left transition-colors"
							@click="handleAddLayerType(option.type)"
						>
							<div class="flex items-center gap-1.5">
								<Icon :name="option.icon" class="h-4 w-4 shrink-0" />
								<span class="text-sm font-medium">{{ option.label }}</span>
							</div>
							<p class="text-muted-foreground text-xs">{{ option.description }}</p>
						</button>
					</div>
				</PopoverContent>
			</Popover>
		</div>
		<CardContent class="flex-1 overflow-hidden p-0">
			<ScrollArea class="h-full">
				<div v-if="layers.length === 0" class="p-8 text-center text-muted-foreground">
					<p class="text-sm">Слоёв пока нет</p>
					<p class="mt-1 text-xs">Нажмите «Добавить», чтобы создать слой</p>
				</div>
				<VueDraggable
					v-if="displayLayers.length > 0"
					v-model="displayLayers"
					:animation="150"
					:filter="'input, button'"
					:prevent-on-filter="false"
					ghost-class="opacity-30"
					class="space-y-0.5 p-1.5"
					@end="handleReorder"
				>
					<div
						v-for="layer in displayLayers"
						:id="`layer-row-${layer.id}`"
						:key="layer.id"
						class="group flex h-8 cursor-pointer items-center gap-2 rounded-md px-1.5 transition-colors"
						:class="[
							isLayerSelected(layer.id)
								? 'bg-emerald-500/10 text-emerald-700 shadow-[inset_2px_0_0_#10b981] hover:bg-emerald-500/[0.13] dark:text-emerald-300'
								: 'text-foreground hover:bg-accent/70',
							{ 'opacity-50': !layer.visible },
						]"
						@click="handleLayerClick(layer.id, $event)"
					>
						<span class="inline-flex size-[22px] flex-none items-center justify-center rounded-md" :class="getLayerTypeMeta(layer.type).chipClass">
							<Icon :name="getLayerTypeMeta(layer.type).icon" class="h-3.5 w-3.5" />
						</span>

						<input
							v-if="renamingLayerId === layer.id"
							id="layer-rename-input"
							v-model="renameDraft"
							type="text"
							class="h-5 min-w-0 flex-1 rounded border border-emerald-500/70 bg-background px-1 text-xs text-foreground outline-none ring-2 ring-emerald-500/15"
							@click.stop
							@dblclick.stop
							@keydown.stop
							@keydown.enter.prevent="commitRename"
							@keydown.esc.prevent="cancelRename"
							@blur="commitRename"
						/>
						<span v-else class="min-w-0 flex-1 select-none truncate text-xs" @dblclick.stop="startRename(layer)">{{ layer.name }}</span>

						<span class="flex flex-none items-center gap-1.5 group-hover:hidden">
							<span class="text-[10px] tabular-nums text-muted-foreground">{{ layer.width }}×{{ layer.height }}</span>
							<Icon v-if="layer.locked" name="lucide:lock" class="h-3 w-3 text-muted-foreground/60" />
							<Icon v-if="!layer.visible" name="lucide:eye-off" class="h-3 w-3 text-muted-foreground" />
						</span>

						<span class="hidden flex-none items-center gap-0.5 group-hover:flex">
							<button type="button" class="inline-flex size-6 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-foreground" :title="layer.visible ? 'Скрыть' : 'Показать'" @click.stop="emit('toggleVisibility', layer.id)">
								<Icon :name="layer.visible ? 'lucide:eye' : 'lucide:eye-off'" class="h-3.5 w-3.5" />
							</button>
							<button type="button" class="inline-flex size-6 items-center justify-center rounded-md transition-colors hover:bg-accent" :class="layer.locked ? 'text-amber-600 hover:text-amber-700 dark:text-amber-400/80 dark:hover:text-amber-300' : 'text-muted-foreground hover:text-foreground'" :title="layer.locked ? 'Разблокировать' : 'Заблокировать'" @click.stop="emit('toggleLock', layer.id)">
								<Icon :name="layer.locked ? 'lucide:lock' : 'lucide:lock-open'" class="h-3.5 w-3.5" />
							</button>
						</span>
					</div>
				</VueDraggable>
			</ScrollArea>
		</CardContent>
	</Card>
</template>
