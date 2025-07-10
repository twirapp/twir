import { createGlobalState } from '@vueuse/core'
import { readonly, ref } from 'vue'
import type { KappagenConfig } from '@twirapp/kappagen/types'

export const useKappagenSettings = createGlobalState(() => {
	const settings = ref<KappagenConfig>(null)

	function setSettings(newSettings: KappagenConfig) {
		settings.value = newSettings
	}

	return {
		settings: readonly(settings),
		setSettings,
	}
})
