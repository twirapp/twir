import { type Ref, computed } from 'vue'

import type { Layer } from '../types'

type UpdateLayer = (updates: Partial<Layer>) => void

export function useLayerProperties(layer: Ref<Layer | null>, updateLayer: UpdateLayer) {
	const fieldId = (name: string) => `layer-props-${layer.value?.id ?? ''}-${name}`
	const update = (updates: Partial<Layer>) => {
		if (layer.value) updateLayer(updates)
	}
	const localName = computed({ get: () => layer.value?.name ?? '', set: (value: string) => update({ name: value }) })
	const localPosX = computed({ get: () => layer.value?.posX ?? 0, set: (value: number) => update({ posX: value }) })
	const localPosY = computed({ get: () => layer.value?.posY ?? 0, set: (value: number) => update({ posY: value }) })
	const localWidth = computed({ get: () => layer.value?.width ?? 0, set: (value: number) => update({ width: value }) })
	const localHeight = computed({ get: () => layer.value?.height ?? 0, set: (value: number) => update({ height: value }) })
	const localRotation = computed({ get: () => layer.value?.rotation ?? 0, set: (value: number) => update({ rotation: value }) })
	const localOpacity = computed({ get: () => (layer.value?.opacity ?? 1) * 100, set: (value: number) => update({ opacity: value / 100 }) })
	const localVisible = computed({ get: () => layer.value?.visible ?? true, set: (value: boolean) => update({ visible: value }) })
	const localLocked = computed({ get: () => layer.value?.locked ?? false, set: (value: boolean) => update({ locked: value }) })
	const localPeriodicallyRefetch = computed({
		get: () => layer.value?.periodicallyRefetchData ?? true,
		set: (value: boolean) => update({ periodicallyRefetchData: value }),
	})
	const localPollInterval = computed({
		get: () => layer.value?.settings.htmlOverlayDataPollSecondsInterval ?? 5,
		set: (value: number) => {
			if (layer.value) update({ settings: { ...layer.value.settings, htmlOverlayDataPollSecondsInterval: value } })
		},
	})

	return {
		fieldId,
		localName,
		localPosX,
		localPosY,
		localWidth,
		localHeight,
		localRotation,
		localOpacity,
		localVisible,
		localLocked,
		localPeriodicallyRefetch,
		localPollInterval,
	}
}
