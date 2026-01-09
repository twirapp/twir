import { createGlobalState } from '@vueuse/core'
import { ref, toRaw } from 'vue'

import { type ChatSettingsWithOptionalId, defaultChatSettings } from './default-settings'

export const useChatOverlayForm = createGlobalState(() => {
	const data = ref<ChatSettingsWithOptionalId>(structuredClone(defaultChatSettings))

	function setData(d: ChatSettingsWithOptionalId) {
		data.value = structuredClone(toRaw(d))
	}

	function reset() {
		data.value = {
			id: data.value.id,
			...structuredClone(defaultChatSettings),
		}
	}

	function getDefaultSettings() {
		return structuredClone(defaultChatSettings)
	}

	return {
		data,
		setData,
		reset,
		getDefaultSettings,
	}
})
