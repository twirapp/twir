import { createGlobalState, useLocalStorage } from '@vueuse/core'

export const useWidgetPreviewMode = createGlobalState(() => {
	const enabled = useLocalStorage('overlay-builder-widgets-preview', false)

	return { enabled }
})
