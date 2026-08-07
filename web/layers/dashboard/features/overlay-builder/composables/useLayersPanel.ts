import { type Ref, nextTick, ref, watch } from 'vue'
import { ChannelOverlayLayerType } from '~/gql/graphql.js'

import type { Layer } from '../types'

import { useLayerRename } from './useLayerRename'

export interface LayerTypeOption {
	readonly type: ChannelOverlayLayerType
	readonly icon: string
	readonly label: string
	readonly description: string
}

const layerTypeOptions: readonly LayerTypeOption[] = [
	{ type: ChannelOverlayLayerType.Image, icon: 'lucide:image', label: 'Картинка', description: 'Изображение из URL' },
	{ type: ChannelOverlayLayerType.Video, icon: 'lucide:video', label: 'Видео', description: 'Видео из URL' },
	{ type: ChannelOverlayLayerType.Youtube, icon: 'simple-icons:youtube', label: 'YouTube', description: 'Видео с YouTube' },
	{ type: ChannelOverlayLayerType.Text, icon: 'lucide:type', label: 'Текст', description: 'Текстовый слой' },
	{ type: ChannelOverlayLayerType.Html, icon: 'lucide:code-xml', label: 'HTML', description: 'HTML, CSS и JavaScript' },
	{ type: ChannelOverlayLayerType.Iframe, icon: 'lucide:panels-top-left', label: 'Виджет', description: 'Встраиваемый URL' },
	{ type: ChannelOverlayLayerType.Emote, icon: 'lucide:smile', label: 'Эмоции', description: 'Один эмоут на слой' },
]

type SelectLayer = (layerId: string, addToSelection: boolean) => void
type AddLayer = (type: ChannelOverlayLayerType) => void
type ReorderLayers = (layers: Layer[]) => void
type UpdateLayer = (layerId: string, updates: Partial<Layer>) => void

export function useLayersPanel(
	layers: Ref<Layer[]>,
	selectedLayerIds: Ref<string[]>,
	selectLayer: SelectLayer,
	addLayer: AddLayer,
	reorderLayers: ReorderLayers,
	updateLayer: UpdateLayer,
) {
	const displayLayers = ref<Layer[]>([])
	const isAddPopoverOpen = ref(false)

	const rename = useLayerRename(layers, updateLayer)

	function handleAddLayerType(type: ChannelOverlayLayerType) {
		isAddPopoverOpen.value = false
		addLayer(type)
	}

	watch(
		layers,
		(newLayers) => {
			displayLayers.value = [...newLayers].reverse()
		},
		{ immediate: true, deep: true }
	)

	watch(selectedLayerIds, (ids) => {
		if (ids.length !== 1) return
		nextTick(() => {
			document.getElementById(`layer-row-${ids[0]}`)?.scrollIntoView({ block: 'nearest' })
		})
	})

	function handleReorder() {
		reorderLayers([...displayLayers.value].reverse())
	}

	function handleLayerClick(layerId: string, event: MouseEvent) {
		selectLayer(layerId, event.ctrlKey || event.metaKey)
	}

	const isLayerSelected = (layerId: string) => selectedLayerIds.value.includes(layerId)

	return {
		displayLayers,
		isAddPopoverOpen,
		layerTypeOptions,
		handleAddLayerType,
		handleReorder,
		handleLayerClick,
		isLayerSelected,
		...rename,
	}
}
