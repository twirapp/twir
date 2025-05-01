import { defineStore } from 'pinia'

import type { PasteBinOutputDto } from '@twir/api/openapi'

export const usePasteStore = defineStore('paste', () => {
	const currentPaste = ref<PasteBinOutputDto>()
	const editableContent = ref<string>()

	function setCurrentPaste(paste: PasteBinOutputDto) {
		currentPaste.value = paste
	}

	function duplicate() {
		if (!currentPaste.value) return

		editableContent.value = currentPaste.value.content
		currentPaste.value = undefined
	}

	return {
		currentPaste: readonly(currentPaste),
		setCurrentPaste,

		editableContent,

		duplicate,
	}
})
