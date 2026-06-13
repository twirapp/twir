import { createGlobalState } from '@vueuse/core'
import { ref, toRaw } from 'vue'

export const useDudesIframe = createGlobalState(() => {
	const dudesIframe = ref<HTMLIFrameElement | null>(null)

	function sendIframeMessage(action: string, data?: any) {
		if (!dudesIframe.value) return
		const payload = JSON.stringify({
			action,
			data: toRaw(data),
		})
		dudesIframe.value.contentWindow?.postMessage(payload)
	}

	return {
		dudesIframe,
		sendIframeMessage,
	}
})
