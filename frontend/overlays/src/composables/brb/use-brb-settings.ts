import { createGlobalState } from '@vueuse/core'
import { ref } from 'vue'

import type { BrbSetSettingsFn } from '@/types.js'

export interface BrbSettings {
	text: string
	late: {
		enabled: boolean
		text: string
		displayBrbTime: boolean
	}
	backgroundColor: string
	fontSize: number
	fontColor: string
	fontFamily: string
}

export const useBrbSettings = createGlobalState(() => {
	const settings = ref<BrbSettings>()

	const setSettings: BrbSetSettingsFn = (newSettings) => {
		settings.value = newSettings
	}

	return {
		settings,
		setSettings,
	}
})
