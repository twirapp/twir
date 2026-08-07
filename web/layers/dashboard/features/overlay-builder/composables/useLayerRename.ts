import { type Ref, nextTick, ref } from 'vue'

import type { Layer } from '../types'

type UpdateLayer = (layerId: string, updates: Partial<Layer>) => void

export function useLayerRename(layers: Ref<Layer[]>, updateLayer: UpdateLayer) {
	const renamingLayerId = ref<string | null>(null)
	const renameDraft = ref('')

	function startRename(layer: Layer) {
		renamingLayerId.value = layer.id
		renameDraft.value = layer.name
		nextTick(() => {
			const input = document.querySelector<HTMLInputElement>('#layer-rename-input')
			input?.focus()
			input?.select()
		})
	}

	function commitRename() {
		if (renamingLayerId.value === null) return
		const layerId = renamingLayerId.value
		renamingLayerId.value = null

		const name = renameDraft.value.trim()
		const layer = layers.value.find((item) => item.id === layerId)
		if (layer && name && name !== layer.name) updateLayer(layerId, { name })
	}

	function cancelRename() {
		renamingLayerId.value = null
	}

	return { renamingLayerId, renameDraft, startRename, commitRename, cancelRename }
}
