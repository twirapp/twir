import { type Ref, computed } from 'vue'

import type { Layer, LayerSettings } from '../types'

type UpdateSettings = (updates: Partial<LayerSettings>) => void

export function useImageLayerEditor(layer: Ref<Layer>, updateSettings: UpdateSettings) {
	const imageUrl = computed({
		get: () => layer.value.settings.imageUrl,
		set: (value: string) => updateSettings({ imageUrl: value }),
	})

	function setPlaceholder() {
		updateSettings({ imageUrl: 'https://via.placeholder.com/300x200' })
	}

	return { imageUrl, setPlaceholder }
}
