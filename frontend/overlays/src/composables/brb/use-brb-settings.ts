import { createGlobalState } from '@vueuse/core'
import { ref } from 'vue'

import type { BrbSetSettingsFn } from '@/types.js'
import type { Settings } from '@twir/api/messages/overlays_be_right_back/overlays_be_right_back'

export const useBrbSettings = createGlobalState(() => {
	const settings = ref<Settings>()

	const setSettings: BrbSetSettingsFn = (newSettings) => {
		settings.value = newSettings
	}

	return {
		settings,
		setSettings,
	}
})
