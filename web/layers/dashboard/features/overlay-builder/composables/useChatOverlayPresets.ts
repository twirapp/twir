import { until } from '@vueuse/core'
import type { AcceptableValue } from 'reka-ui'

import { useChatOverlayApi } from '~~/layers/dashboard/api/overlays/chat.ts'
import { useChatOverlayForm } from '~~/layers/dashboard/pages/dashboard/overlays/chat/components/form.ts'

type SelectPreset = (id: string) => void

export function useChatOverlayPresets(onSelectPreset: SelectPreset) {
	const chatOverlaysManager = useChatOverlayApi()
	const creator = chatOverlaysManager.useOverlayCreate()
	const { data: chatOverlaysData, fetching: fetchingOverlays } = chatOverlaysManager.useOverlaysQuery()
	const { setData, getDefaultSettings } = useChatOverlayForm()

	const selectedPresetId = ref<string>()
	const presets = computed(() => chatOverlaysData.value?.chatOverlays ?? [])

	watch(
		() => chatOverlaysData.value?.chatOverlays,
		(overlays) => {
			if (!overlays?.length) {
				selectedPresetId.value = undefined
				return
			}

			if (!overlays.some((overlay) => overlay.id === selectedPresetId.value)) {
				selectedPresetId.value = overlays[0]?.id ?? undefined
			}
		},
		{ immediate: true }
	)

	watch(selectedPresetId, (id) => {
		const preset = presets.value.find((overlay) => overlay.id === id)
		if (preset) setData(preset)
		if (id) onSelectPreset(id)
	})

	function handlePresetChange(id: AcceptableValue) {
		if (typeof id !== 'string') return

		selectedPresetId.value = id
		const preset = presets.value.find((overlay) => overlay.id === id)
		if (preset) setData(preset)
	}

	async function createPreset() {
		const previousLength = presets.value.length

		await creator.executeMutation({ input: getDefaultSettings() })
		await until(() => chatOverlaysData.value?.chatOverlays).changed()

		const overlays = chatOverlaysData.value?.chatOverlays ?? []
		if (overlays.length > previousLength) {
			selectedPresetId.value = overlays[overlays.length - 1]?.id ?? undefined
		}
	}

	return {
		fetchingOverlays,
		presets,
		selectedPresetId,
		handlePresetChange,
		createPreset,
	}
}
