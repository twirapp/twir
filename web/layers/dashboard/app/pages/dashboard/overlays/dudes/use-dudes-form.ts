import { ref, toRaw } from 'vue'

import { type DudesSettingsWithOptionalId, defaultDudesSettings } from './dudes-settings.js'

const data = ref<DudesSettingsWithOptionalId>(structuredClone(defaultDudesSettings))

export function useDudesForm() {
	function setData(d: DudesSettingsWithOptionalId) {
		data.value = structuredClone(toRaw(d))
	}

	function reset() {
		data.value = {
			...getDefaultSettings(),
			id: data.value.id,
		}
	}

	function getDefaultSettings() {
		return structuredClone(defaultDudesSettings)
	}

	return {
		data,
		setData,
		reset,
		getDefaultSettings,
	}
}
