import { defineStore } from 'pinia'

import type { PasteBinOutputDto } from '@twir/api/openapi'

export const usePasteStore = defineStore('paste', () => {
	const currentPaste = ref<PasteBinOutputDto>()
	const editableContent = ref<string>()

	function setCurrentPaste(paste?: PasteBinOutputDto) {
		currentPaste.value = paste
	}

	return {
		currentPaste: readonly(currentPaste),
		setCurrentPaste,

		editableContent,
	}
})

if (import.meta.hot) {
	import.meta.hot.accept(acceptHMRUpdate(usePasteStore, import.meta.hot))
}
