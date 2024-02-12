import { Settings } from '@twir/api/messages/overlays_now_playing/overlays_now_playing';
import { defineStore } from 'pinia';
import { ref, toRaw } from 'vue';

export const useNowPlayingForm = defineStore('now-playing-form', () => {
	const data = ref<Settings>();

	function $setData(d: Settings) {
		data.value = structuredClone(toRaw(d));
	}

	function $getDefaultSettings() {
		return structuredClone(Settings);
	}

	return {
		data,
		$setData,
		$getDefaultSettings,
	};
});
