import { ref, toRaw } from 'vue';

import { defaultDudesSettings, type DudesSettingsWithOptionalId } from './dudes-settings.js';

const data = ref<DudesSettingsWithOptionalId>(structuredClone(defaultDudesSettings));

export const useDudesForm = () => {
	function $setData(d: DudesSettingsWithOptionalId) {
		data.value = structuredClone(toRaw(d));
	}

	function $reset() {
		data.value = structuredClone(defaultDudesSettings);
	}

	function $getDefaultSettings() {
		return structuredClone(defaultDudesSettings);
	}

	return {
		data,
		$setData,
		$reset,
		$getDefaultSettings,
	};
};
