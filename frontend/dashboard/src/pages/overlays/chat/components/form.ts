import { ref, toRaw } from 'vue';

import { defaultChatSettings, type ChatSettingsWithOptionalId } from './default-settings';

const data = ref<ChatSettingsWithOptionalId>(structuredClone(defaultChatSettings));

export const useChatOverlayForm = () => {
	function $setData(d: ChatSettingsWithOptionalId) {
		data.value = structuredClone(toRaw(d));
	}

	function $reset() {
		data.value = {
			id: data.value.id,
			...structuredClone(defaultChatSettings),
		};
	}

	function $getDefaultSettings() {
		return structuredClone(defaultChatSettings);
	}

	return {
		data,
		$setData,
		$reset,
		$getDefaultSettings,
	};
};
