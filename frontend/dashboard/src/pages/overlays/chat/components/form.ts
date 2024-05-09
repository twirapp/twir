import { ref, toRaw } from 'vue'

import { type ChatSettingsWithOptionalId, defaultChatSettings } from './default-settings'

const data = ref<ChatSettingsWithOptionalId>(structuredClone(defaultChatSettings))

export function useChatOverlayForm() {
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
}
